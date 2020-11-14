package heidou

import (
	"fmt"
	"html/template"
)

var DefaultMethods = []string{"list", "create", "get", "update", "delete", "bulkGet", "bulkDelete"}

type ErrorCode string
type Method string

// MetaTypes mappings for sql types to json, go etc
type MetaTypes struct {
	MetaTypes []MetaType `yaml:"mappings"`
}

// MetaType mapping
type MetaType struct {
	// SQLType sql type reported from db
	SQLType string `yaml:"sqlType"`

	// GoType mapped go type
	GoType string `yaml:"goType"`

	// GqlType mapped graphql type
	GqlType string `yaml:"gqlType"`

	// JSONType mapped json type
	JSONType string `yaml:"jsonType"`

	// ProtobufType mapped protobuf type
	ProtobufType string `yaml:"protobufType"`

	// GureguType mapped go type using Guregu
	GureguType string `yaml:"gureguType"`

	// GoNullableType mapped go type using nullable
	GoNullableType string `yaml:"goNullableType"`

	// SwaggerType mapped type
	SwaggerType string `yaml:"swaggerType"`
}

func (m *MetaType) String() interface{} {
	return fmt.Sprintf("SQLType: %-15s  GoType: %-15s GqlType: %-15s GureguType: %-15s GoNullableType: %-15s JSONType: %-15s ProtobufType: %-15s",
		m.SQLType, m.GqlType,
		m.GoType, m.GureguType, m.GoNullableType,
		m.JSONType, m.ProtobufType)
}

type Column struct {
	Name     string
	Type     string
	DataType string
	Comment  string
	Key      string
}
type MetaTable struct {
	Name    string
	Columns []*Column
}

type Table struct {
	Name        string `yaml:"name"`
	PkgPath     string `yaml:"pkgPath"`
	Description string `yaml:"description"`
	IsSkip      bool   `yaml:"isSkip"`

	Fields     []*Field    `yaml:"fields"`
	ErrorCodes []ErrorCode `yaml:"errorCodes"`
	Methods    []string    `yaml:"methods"`

	PrimaryKeyField *Field
	Filterable      bool
	Sortable        bool
	HasErrorCode    bool
	HasTimeField    bool
	IsStringsJoin   bool

	NameSnake            string
	NameSnakePlural      string
	NameKebab            string
	NameKebabPlural      string
	NameCamel            string
	NameCamelPlural      string
	NameLowerCamel       string
	NameLowerCamelPlural string
}

func (t *Table) genName() {
	name := rules.Singularize(t.Name)
	t.NameSnake = snake(name)
	t.NameSnakePlural = plural(t.NameSnake)
	t.NameKebab = snake(name)
	t.NameKebabPlural = plural(t.NameKebab)
	t.NameCamel = pascal(name)
	t.NameCamelPlural = plural(t.NameCamel)
	t.NameLowerCamel = camel(name)
	t.NameLowerCamelPlural = plural(t.NameLowerCamel)

}

// 查找配置中的表信息
func getTableInCfg(tables []*Table, name string) *Table {
	for _, table := range tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

// 合并数据库定义和项目配置中的表信息, 构建为新表结构
// 返回 nil 表示 不生成此表信息
func mergeTable(metaTable *MetaTable, tableInCfg *Table, metaTypes map[string]MetaType) *Table {
	table := &Table{
		Name: metaTable.Name,
	}

	// 生成数据库定义的字段
	for _, column := range metaTable.Columns {
		field := shiftMetaField(column, metaTypes)

		var fieldInCfg *Field
		if tableInCfg != nil {
			fieldInCfg = findField(tableInCfg.Fields, column.Name)
		}

		field = mergeField(field, fieldInCfg)

		if field.IsPrimaryKey {
			table.PrimaryKeyField = field
		}

		if field.IsFilterable {
			table.Filterable = true
			for _, op := range field.Operations {
				if op == "In" {
					table.IsStringsJoin = true
				}
			}
		}

		if field.IsSortable {
			table.Sortable = true
		}

		if field.MetaType.GoType == "time.Time" {
			table.HasTimeField = true
		}

		if table.Name == "product" {
			fmt.Println("field:", field)
		}
		table.Fields = append(table.Fields, field)
	}
	table.Methods = DefaultMethods

	if tableInCfg != nil {
		if len(tableInCfg.ErrorCodes) > 0 {
			table.HasErrorCode = true
			table.ErrorCodes = append(table.ErrorCodes, tableInCfg.ErrorCodes...)
		}

		// 配置了跳过，则不生成此表的信息
		table.IsSkip = tableInCfg.IsSkip

		if len(tableInCfg.Methods) > 0 {
			table.Methods = tableInCfg.Methods
		}

		// 生成配置中关联类型的字段
		for _, field := range tableInCfg.Fields {
			field.genName()

			if field.JoinType == "" {
				continue
			}

			if field.IsFilterable {
				table.Filterable = true
			}

			if field.IsSortable {
				table.Sortable = true
			}

			tags := `json:"` + field.NameLowerCamel + `" gorm:"` + field.NameLowerCamel

			if field.JoinTableName != "" {
				tags += ";many2many:" + field.JoinTableName
			}
			if field.ForeignKey != "" {
				tags += ";foreignKey:" + field.ForeignKey
			}
			if field.References != "" {
				tags += ";references:" + field.References
			}
			if field.JoinForeignKey != "" {
				tags += ";joinForeignKey:" + field.JoinForeignKey
			}
			if field.JoinReferences != "" {
				tags += ";joinReferences:" + field.JoinReferences
			}
			tags += `"`

			fmt.Println("tags:", tags)
			field.TagsHTML = template.HTML(tags)

			table.Fields = append(table.Fields, field)
		}
	}

	if table.Name == "product" {

		fmt.Println("Table:", table.Name)
		for _, field := range table.Fields {
			fmt.Println("Field:", field, field.MetaType)
		}
	}
	return table
}
