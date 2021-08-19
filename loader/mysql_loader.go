package loader

import (
	"database/sql"
	"errors"
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

func (msl *MysqlSchemaLoader) loadIndexes(db *sql.DB, table *heidou.MetaTable) error {

	rawSql := "SELECT `INDEX_NAME`, `COLUMN_NAME`,`NON_UNIQUE`, `SEQ_IN_INDEX` FROM `STATISTICS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? ORDER BY table_name ASC, SEQ_IN_INDEX ASC"
	fmt.Printf("load index:%s\n", rawSql)
	rows, err := db.Query(rawSql, msl.SchemaName, table.Name)
	if err != nil {
		return err
	}
	for rows.Next() {
		var indexName, columnName string
		var nonUnique, seq int
		if rows.Scan(&indexName, &columnName, &nonUnique, &seq) == nil {
			index := &heidou.Index{
				Name:       indexName,
				ColumnName: columnName,
				Unique:     nonUnique == 0,
				Seq:        seq,
			}

			table.Indexes = append(table.Indexes, index)

		} else {
			logrus.Errorf("get %s database %s table's index info failed", msl.SchemaName, table.Name)
		}
	}
	return nil
}

func (msl *MysqlSchemaLoader) LoadMetaTable() ([]*heidou.MetaTable, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/INFORMATION_SCHEMA?charset=%s", msl.User, msl.Password, msl.Host, msl.Port, msl.Charset)
	db, err := sql.Open(DialectMysql, dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rawSql := "SELECT COUNT(*) FROM `COLUMNS` WHERE `TABLE_SCHEMA` = ? "
	row := db.QueryRow(rawSql, msl.SchemaName)
	count := 0
	err = row.Scan(&count)
	if err != nil {
		return nil, err
	}

	rawSql = "SELECT `TABLE_NAME`, `COLUMN_NAME`,`DATA_TYPE`,`COLUMN_TYPE`,`COLUMN_COMMENT`, `EXTRA`, `IS_NULLABLE` FROM `COLUMNS` WHERE `TABLE_SCHEMA` = ? ORDER BY table_name ASC, ordinal_position ASC"
	rows, err := db.Query(rawSql, msl.SchemaName)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New(fmt.Sprintf("can't find table %s's meta data", msl.SchemaName))
	}

	tables := make([]*heidou.MetaTable, 0)
	tablesIndex := make(map[string]*heidou.MetaTable)

	for rows.Next() {
		var tableName, columnName, dataType, columnType, columnComment, extra, nullable string

		if rows.Scan(&tableName, &columnName, &dataType, &columnType, &columnComment, &extra, &nullable) == nil {
			c := &heidou.Column{
				Name: columnName,
				// DataType: dataType,
				Type:    dataType,
				Comment: columnComment,
				// Extra:   extra,
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

	for _, table := range tables {
		err := msl.loadIndexes(db, table)
		if err != nil {
			return nil, err
		}
		for _, column := range table.Columns {
			fmt.Printf("%s %#v\n", table.Name, column)
			for _, index := range table.Indexes {
				fmt.Printf("%s %s\n", index.ColumnName, column.Name)
				if index.ColumnName == column.Name {
					fmt.Printf("%#v\n", index)
					column.IsIndex = true

					if index.Unique {
						column.IsUnique = true
					}
					if index.Name == "PRIMARY" {
						column.IsPrimaryKey = true
					}
				}

			}

		}
	}
	return tables, nil
}
