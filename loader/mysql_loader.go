package loader

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/ychengcloud/heidou"
)

const (
	DialectMysql = "mysql"
)

type MysqlSchemaLoaderOptions struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Charset  string `yaml:"charset"`
}

type MysqlSchemaLoader struct {
	*MysqlSchemaLoaderOptions
	SchemaName string
}

func NewMysqlSchemaLoader(o *MysqlSchemaLoaderOptions, schemaName string) *MysqlSchemaLoader {
	return &MysqlSchemaLoader{
		MysqlSchemaLoaderOptions: o,
		SchemaName:               schemaName,
	}
}

func (msl *MysqlSchemaLoader) LoadMetaTable() ([]*heidou.MetaTable, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/INFORMATION_SCHEMA?charset=%s", msl.User, msl.Password, msl.Host, msl.Port, msl.Charset)
	db, err := sql.Open(DialectMysql, dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rawSql := "SELECT `TABLE_NAME`, `COLUMN_NAME`,`DATA_TYPE`,`COLUMN_TYPE`,`COLUMN_COMMENT`,`COLUMN_KEY`, `EXTRA` FROM `COLUMNS` WHERE `TABLE_SCHEMA` = ?"
	rows, err := db.Query(rawSql, msl.SchemaName)
	if err != nil {
		return nil, err
	}

	tables := make([]*heidou.MetaTable, 0)
	tablesIndex := make(map[string]*heidou.MetaTable)

	for rows.Next() {
		var tableName, columnName, dataType, columnType, columnComment, columnKey, extra string

		if rows.Scan(&tableName, &columnName, &dataType, &columnType, &columnComment, &columnKey, &extra) == nil {
			c := &heidou.Column{
				Name:     columnName,
				DataType: dataType,
				Type:     columnType,
				Comment:  columnComment,
				Key:      columnKey,
				Extra:    extra,
			}

			if strings.ToUpper(columnKey) == "PRI" {
				c.IsPrimaryKey = true
			}

			if strings.Contains(extra, "auto_increment") {
				c.IsAutoIncrement = true
			}

			table, ok := tablesIndex[tableName]
			if !ok {
				table = &heidou.MetaTable{
					Name:    tableName,
					Columns: make([]*heidou.Column, 0),
				}
				tablesIndex[tableName] = table

				tables = append(tables, table)

			}
			table.Columns = append(table.Columns, c)

		} else {
			logrus.Errorf("get %s database column info failed", msl.SchemaName)
		}
	}
	return tables, nil
}
