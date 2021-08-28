package heidou

import (
	"html/template"
	"strconv"
	"strings"

	"github.com/jinzhu/inflection"
)

const (
	OpEq         = "Eq"
	OpIn         = "In"
	OpGt         = "Gt"
	OpGte        = "Gte"
	OpLt         = "Lt"
	OpLte        = "Lte"
	OpContains   = "Contains"
	OpStartsWith = "StartsWith"
	OpEndsWith   = "EndsWith"
	OpAnd        = "AND"
	OpOr         = "OR"
	OpNot        = "NOT"
)

type JoinType string

const (
	JoinTypeNone       JoinType = "None"
	JoinTypeBelongTo   JoinType = "BelongTo"
	JoinTypeHasOne     JoinType = "HasOne"
	JoinTypeHasMany    JoinType = "HasMany"
	JoinTypeManyToMany JoinType = "ManyToMany"
)

type Field struct {
	//从配置文件读取的数据
	Name         string   `mapstructure:"name" yaml:"name"`
	Alias        string   `mapstructure:"alias" yaml:"alias"`
	Description  string   `mapstructure:"description" yaml:"description"`
	Tags         string   `mapstructure:"tags" yaml:"tags"`
	IsSkip       bool     `mapstructure:"isSkip" yaml:"isSkip"`
	IsRequired   bool     `mapstructure:"isRequired" yaml:"isRequired"`
	IsSortable   bool     `mapstructure:"isSortable" yaml:"isSortable"`
	IsFilterable bool     `mapstructure:"isFilterable" yaml:"isFilterable"`
	JoinType     JoinType `mapstructure:"joinType" yaml:"joinType"`

	//Default: {field_name}
	RefTableName string `mapstructure:"refTableName" yaml:"refTableName"`

	//Default: {table}_{ref_table}
	JoinTableName string `mapstructure:"joinTableName" yaml:"joinTableName"`

	//Default: {ref_table}_id
	ForeignKey string `mapstructure:"foreignKey" yaml:"foreignKey"`

	//Default: id
	References string `mapstructure:"references" yaml:"references"`

	//Default: {table}_id
	JoinForeignKey string `mapstructure:"joinForeignKey" yaml:"joinForeignKey"`

	//Default: {ref_table}_id
	JoinReferences string   `mapstructure:"joinReferences" yaml:"joinReferences"`
	Operations     []string `mapstructure:"operations" yaml:"operations"`

	//生成的数据
	IsPrimaryKey    bool
	IsForeignKey    bool
	IsAutoIncrement bool
	IsUnique        bool
	IsIndex         bool
	MetaType        MetaType
	TagsHTML        template.HTML

	//字段所在表
	Table *Table

	//引用表
	RefTable *Table

	//如果是M2M关联字段，字段对应的联接表
	JoinTable *Table

	MaxLength int

	NameSnake            string
	NameSnakePlural      string
	NameKebab            string
	NameKebabPlural      string
	NameCamel            string
	NameCamelPlural      string
	NameLowerCamel       string
	NameLowerCamelPlural string
}

func (f *Field) genName() {
	name := f.Name
	f.NameSnake = snake(name)
	f.NameSnakePlural = inflection.Plural(f.NameSnake)
	f.NameKebab = snake(name)
	f.NameKebabPlural = inflection.Plural(f.NameKebab)
	f.NameCamel = pascal(name)
	f.NameCamelPlural = inflection.Plural(f.NameCamel)
	f.NameLowerCamel = camel(name)
	f.NameLowerCamelPlural = inflection.Plural(f.NameLowerCamel)
}

func findField(fields []*Field, name string) *Field {
	for _, v := range fields {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// parseColumnType parse sql type and return raw type and length
func parseColumnLength(columnType string) (maxLength int) {

	columnType = strings.ToLower(columnType)
	maxLength = -1
	idx1 := strings.Index(columnType, "(")
	idx2 := strings.Index(columnType, ")")

	if idx1 > -1 && idx2 > -1 {
		sizeStr := columnType[idx1+1 : idx2]
		i, err := strconv.Atoi(sizeStr)
		if err == nil {
			maxLength = i
		}
	}

	return maxLength
}

func shiftMetaField(column *Column, metaTypes map[string]MetaType) *Field {

	// isUnsigned := strings.Contains(column.Type, " unsigned") || strings.Contains(column.Type, " UNSIGNED")

	columnType := strings.ToLower(column.Type)
	// if isUnsigned {
	// 	columnType = "u" + columnType
	// }

	// maxLength := parseColumnLength(column.Type)
	//TODO
	maxLength := 0
	field := &Field{
		Name:            column.Name,
		Description:     column.Comment,
		MetaType:        metaTypes[columnType], //meta
		MaxLength:       maxLength,
		JoinType:        JoinTypeNone,
		IsPrimaryKey:    column.IsPrimaryKey,
		IsAutoIncrement: column.IsAutoIncrement,
		IsUnique:        column.IsUnique,
		IsIndex:         column.IsIndex,
	}

	field.genName()

	return field
}

func mergeField(field *Field, fieldInCfg *Field) *Field {

	if fieldInCfg == nil {
		return field
	}

	field.Alias = fieldInCfg.Alias
	field.Description = fieldInCfg.Description
	field.Tags = fieldInCfg.Tags

	if fieldInCfg.JoinType == "" {
		fieldInCfg.JoinType = JoinTypeNone
	} else {
		field.JoinType = fieldInCfg.JoinType
	}

	field.IsSkip = fieldInCfg.IsSkip
	field.IsRequired = fieldInCfg.IsRequired
	field.IsPrimaryKey = fieldInCfg.IsPrimaryKey
	field.IsSortable = fieldInCfg.IsSortable
	field.IsFilterable = fieldInCfg.IsFilterable
	field.Operations = fieldInCfg.Operations
	field.ForeignKey = fieldInCfg.ForeignKey
	field.References = fieldInCfg.References
	field.JoinForeignKey = fieldInCfg.JoinForeignKey
	field.JoinReferences = fieldInCfg.JoinReferences

	return field
}

func (f *Field) HandleAssociation() {
	if f.JoinType == JoinTypeNone {
		return
	}

	if f.RefTableName == "" {
		f.RefTableName = f.Name
	}
	var defaultForeignKey, defaultReferences, defaultJoinForeignKey, defaultJoinReferences string
	switch f.JoinType {
	case JoinTypeBelongTo:
		defaultForeignKey = f.RefTableName + "_id"
		defaultReferences = "id"
	case JoinTypeHasOne:
		defaultForeignKey = f.Table.Name + "_id"
		defaultReferences = f.RefTableName + "_id"
	case JoinTypeHasMany:
		defaultForeignKey = f.Table.Name + "_id"
		defaultReferences = "id"
	case JoinTypeManyToMany:
		defaultForeignKey = "id"
		defaultReferences = "id"
		defaultJoinForeignKey = f.Table.Name + "_id"
		defaultJoinReferences = f.RefTableName + "_id"
	}

	if f.ForeignKey == "" {
		f.ForeignKey = defaultForeignKey
	}
	if f.References == "" {
		f.References = defaultReferences
	}

	if f.JoinType == JoinTypeManyToMany {
		if f.JoinTableName == "" {
			f.JoinTableName = f.Table.Name + "_" + f.RefTableName
		}

		if f.JoinForeignKey == "" {
			f.JoinForeignKey = defaultJoinForeignKey
		}
		if f.JoinReferences == "" {
			f.JoinReferences = defaultJoinReferences
		}
	}

}

func (f *Field) handleOperations() {
	if f.IsIndex || f.IsFilterable {
		if len(f.Operations) == 0 {
			f.Operations = []string{OpEq}
		}
	}
}
