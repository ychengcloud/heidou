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
					Type:            "bigint",
					Comment:         "id",
					IsPrimaryKey:    true,
					IsAutoIncrement: true,
				},
				{
					Name:    "name",
					Type:    "varchar(32)",
					Comment: "name",
				},
				{
					Name:    "age",
					Type:    "int",
					Comment: "age",
				},
			},
		},
		{
			Name: "pet",
			Columns: []*heidou.Column{
				{
					Name:            "id",
					Type:            "bigint",
					Comment:         "id",
					IsPrimaryKey:    true,
					IsAutoIncrement: true,
				},
				{
					Name:    "owner_id",
					Type:    "bigint",
					Comment: "owner_id",
				},
				{
					Name:    "name",
					Type:    "varchar",
					Comment: "name",
				},
			},
		},
		{
			Name: "group",
			Columns: []*heidou.Column{
				{
					Name:            "id",
					Type:            "bigint",
					Comment:         "id",
					IsPrimaryKey:    true,
					IsAutoIncrement: true,
				},
				{
					Name:    "admin_id",
					Type:    "bigint",
					Comment: "admin_id",
				},
				{
					Name:    "name",
					Type:    "varchar(32)",
					Comment: "name",
				},
			},
		},
	}

	return tables, nil
}
