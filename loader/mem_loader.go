package loader

import (
	"github.com/ychengcloud/heidou"
)

type MemoryLoader struct {
}

func NewMemoryLoader() *MemoryLoader {
	return &MemoryLoader{}
}

func (ml *MemoryLoader) LoadMetaTable() ([]*heidou.MetaTable, error) {

	tables := []*heidou.MetaTable{
		{
			Name: "user",
			Columns: []*heidou.Column{
				{
					Name:            "id",
					DataType:        "bigint",
					Type:            "bigint unsigned",
					Comment:         "id",
					Key:             "PRI",
					Extra:           "auto_increment",
					IsPrimaryKey:    true,
					IsAutoIncrement: true,
				},
				{
					Name:     "name",
					DataType: "varchar",
					Type:     "varchar(32)",
					Comment:  "name",
					Key:      "",
				},
				{
					Name:     "age",
					DataType: "int",
					Type:     "int",
					Comment:  "age",
					Key:      "",
				},
			},
		},
		{
			Name: "pet",
			Columns: []*heidou.Column{
				{
					Name:            "id",
					DataType:        "bigint",
					Type:            "bigint unsigned",
					Comment:         "id",
					Key:             "PRI",
					Extra:           "auto_increment",
					IsPrimaryKey:    true,
					IsAutoIncrement: true,
				},
				{
					Name:     "owner_id",
					DataType: "bigint",
					Type:     "bigint unsigned",
					Comment:  "owner_id",
					Key:      "",
				},
				{
					Name:     "name",
					DataType: "varchar",
					Type:     "varchar(32)",
					Comment:  "name",
					Key:      "",
				},
			},
		},
		{
			Name: "group",
			Columns: []*heidou.Column{
				{
					Name:            "id",
					DataType:        "bigint",
					Type:            "bigint unsigned",
					Comment:         "id",
					Key:             "PRI",
					Extra:           "auto_increment",
					IsPrimaryKey:    true,
					IsAutoIncrement: true,
				},
				{
					Name:     "admin_id",
					DataType: "bigint",
					Type:     "bigint unsigned",
					Comment:  "admin_id",
					Key:      "",
				},
				{
					Name:     "name",
					DataType: "varchar",
					Type:     "varchar(32)",
					Comment:  "name",
					Key:      "",
				},
			},
		},
	}

	return tables, nil
}
