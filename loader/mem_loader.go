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
			Name: "post",
			Columns: []*heidou.Column{
				{
					Name:     "id",
					DataType: "bigint",
					Type:     "bigint unsigned",
					Comment:  "id",
					Key:      "PRI",
				},
				{
					Name:     "title",
					DataType: "varchar",
					Type:     "varchar(32)",
					Comment:  "title",
					Key:      "",
				},
				{
					Name:     "content",
					DataType: "text",
					Type:     "text",
					Comment:  "content",
					Key:      "",
				},
			},
		},
		{
			Name: "comment",
			Columns: []*heidou.Column{
				{
					Name:     "id",
					DataType: "bigint",
					Type:     "bigint unsigned",
					Comment:  "id",
					Key:      "PRI",
				},
				{
					Name:     "content",
					DataType: "text",
					Type:     "text",
					Comment:  "content",
					Key:      "",
				},
			},
		},
	}

	return tables, nil
}
