package heidou

import (
	"fmt"

	"github.com/jinzhu/inflection"
)

var DefaultMethods = []string{"list", "create", "get", "update", "delete", "batchCreate", "batchGet", "batchUpdate", "batchDelete"}

type ErrorCode string
type Method string

// MetaTypes mappings for sql types to json, go etc
type MetaTypeMapping struct {
	MetaTypeInfos []MetaTypeInfo `yaml:"mappings"`
}

type MetaTypeInfo struct {
	Dialect   string     `yaml:"dialect"`
	MetaTypes []MetaType `yaml:"types"`
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

	// tsType mapped typescript type
	TSType string `yaml:"tsType"`

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
	Name string
	Type string
	// DataType string
	Comment string
	// Key      string
	// Extra    string

	IsPrimaryKey    bool `yaml:"isPrimaryKey"`
	IsAutoIncrement bool `yaml:"isAutoIncrement"`
	IsUnique        bool `yaml:"isUnique"`
	IsIndex         bool `yaml:"isIndex"`
}
type MetaTable struct {
	Name    string
	Columns []*Column
	Indexes []*Index
}

type Index struct {
	//Name 如果是主键，mysql中始终为 PRIMARY
	Name       string
	ColumnName string
	Unique     bool
	Seq        int
}

type BackReferenceInfo struct {
	Name                 string
	NameCamel            string
	NameCamelPlural      string
	NameLowerCamelPlural string
	JoinTableName        string
}
type Table struct {
	//从配置文件读取的数据
	Name string `mapstructure:"name" yaml:"name"`

	//别名 配置后,生成时不再处理单复数形式,以Alias为准
	Alias string `mapstructure:"alias" yaml:"alias"`

	Description string `mapstructure:"description" yaml:"description"`

	//用于扩展表定义
	Annotations interface{} `mapstructure:"annotations" yaml:"annotations"`

	IsSkip bool `mapstructure:"isSkip" yaml:"isSkip"`

	//是否连接表
	IsJoinTable bool        `mapstructure:"isJoinTable" yaml:"isJoinTable"`
	TypeName    string      `mapstructure:"typeName" yaml:"typeName"` //表类型，可根据此配置选择不同的模板类型
	Extra       interface{} `mapstructure:"extra" yaml:"extra"`

	Fields     []*Field    `mapstructure:"fields" yaml:"fields"`
	ErrorCodes []ErrorCode `mapstructure:"errorCodes" yaml:"errorCodes"`
	Methods    []string    `mapstructure:"methods" yaml:"methods"`

	//IsJoinTable为true时以下字段有效
	JoinFromName   string
	JoinFromFKName string
	JoinToName     string
	JoinToFKName   string

	//生成的数据
	PrimaryKeyField *Field
	//反向引用字段集
	BackReferenceInfos []*BackReferenceInfo
	Filterable         bool
	Sortable           bool
	HasErrorCode       bool
	HasTimeField       bool
	HasJoinField       bool
	IsImportStrings    bool

	NameSnake            string
	NameSnakePlural      string
	NameKebab            string
	NameKebabPlural      string
	NameCamel            string
	NameCamelPlural      string
	NameLowerCamel       string
	NameLowerCamelPlural string
}

type JoinTable struct {
	Name  string
	Table *Table
	//ForeignKey Join ForeignKey field
	ForeignKey *Field `mapstructure:"foreignKey" yaml:"foreignKey"`

	//References Join References field
	References *Field `mapstructure:"references" yaml:"references"`
}

func (t *Table) genName() {
	name := inflection.Singular(t.Name)
	if t.Alias != "" {
		name = t.Alias
	}
	t.NameSnake = snake(name)
	t.NameSnakePlural = inflection.Plural(t.NameSnake)
	t.NameKebab = snake(name)
	t.NameKebabPlural = inflection.Plural(t.NameKebab)
	t.NameCamel = pascal(name)
	t.NameCamelPlural = inflection.Plural(t.NameCamel)
	t.NameLowerCamel = camel(name)
	t.NameLowerCamelPlural = inflection.Plural(t.NameLowerCamel)
}

// 合并数据库定义和项目配置中的表信息, 构建为新表结构
// 返回 nil 表示 不生成此表信息
func MergeTable(metaTable *MetaTable, tableInCfg *Table, metaTypes map[string]MetaType) *Table {
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
		if field.IsSkip {
			continue
		}
		field.handleOperations()
		field.Table = table
		table.handleFlags(field, metaTypes)
		table.Fields = append(table.Fields, field)
	}
	table.Methods = DefaultMethods

	table.handleCfgInfo(tableInCfg)
	return table
}

func (t *Table) handleFlags(field *Field, metaTypes map[string]MetaType) {
	if field.IsPrimaryKey {
		t.PrimaryKeyField = field
	}

	if field.IsFilterable {
		t.Filterable = true
		for _, op := range field.Operations {
			if op == "In" {
				t.IsImportStrings = true
			}
		}
	}

	if field.IsSortable {
		t.Sortable = true
	}

	if field.MetaType.GoType == "time.Time" {
		t.HasTimeField = true
	}
}

func (t *Table) handleCfgInfo(tableInCfg *Table) {
	// 合并配置文件中的表信息
	if tableInCfg == nil {
		return
	}

	t.Alias = tableInCfg.Alias
	t.TypeName = tableInCfg.TypeName
	t.Annotations = tableInCfg.Annotations

	if len(tableInCfg.ErrorCodes) > 0 {
		t.HasErrorCode = true
		t.ErrorCodes = append(t.ErrorCodes, tableInCfg.ErrorCodes...)
	}

	// 配置了跳过，则不生成此表的信息
	t.IsSkip = tableInCfg.IsSkip

	if len(tableInCfg.Methods) > 0 {
		t.Methods = tableInCfg.Methods
	}

	// 生成配置中关联类型的字段
	for _, field := range tableInCfg.Fields {
		if field.JoinType == JoinTypeNone {
			continue
		}
		field.Table = t

		field.genName()
		field.HandleAssociation()

		t.Fields = append(t.Fields, field)
	}
}

func (t *MetaTable) hasCompositeKeys() bool {
	indexMap := make(map[string]*Index, 0)
	for _, index := range t.Indexes {
		//如果map中已经存在，表明是联合索引
		if _, ok := indexMap[index.Name]; ok {
			return true
		}
		indexMap[index.Name] = index
	}

	return false
}
