package loader

import (
	"database/sql"
	"fmt"

	"github.com/decker502/heidou"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const (
	DialectMysql = "mysql"
)

type Options struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Charset  string `yaml:"charset"`
}

type MysqlSchemaLoader struct {
	*Options
	SchemaName string
}

func NewMysqlSchemaLoader(o *Options, schemaName string) *MysqlSchemaLoader {
	return &MysqlSchemaLoader{
		Options:    o,
		SchemaName: schemaName,
	}
}

func (msl *MysqlSchemaLoader) LoadMetaTable() ([]*heidou.MetaTable, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/INFORMATION_SCHEMA?charset=%s", msl.User, msl.Password, msl.Host, msl.Port, msl.Charset)
	db, err := sql.Open(DialectMysql, dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rawSql := "SELECT `TABLE_NAME`, `COLUMN_NAME`,`DATA_TYPE`,`COLUMN_TYPE`,`COLUMN_COMMENT`,`COLUMN_KEY` FROM `COLUMNS` WHERE `TABLE_SCHEMA` = ?"
	rows, err := db.Query(rawSql, msl.SchemaName)
	if err != nil {
		return nil, err
	}

	tables := make([]*heidou.MetaTable, 0)
	tablesIndex := make(map[string]*heidou.MetaTable)

	for rows.Next() {
		var tableName, columnName, dataType, columnType, columnComment, columnKey string

		if rows.Scan(&tableName, &columnName, &dataType, &columnType, &columnComment, &columnKey) == nil {
			c := &heidou.Column{
				Name:     columnName,
				DataType: dataType,
				Type:     columnType,
				Comment:  columnComment,
				Key:      columnKey,
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
