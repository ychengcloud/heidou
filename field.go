package heidou

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
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
	JoinTypeNone       = ""
	JoinTypeBelongTo   = "BelongTo"
	JoinTypeHasOne     = "HasOne"
	JoinTypeHasMany    = "HasMany"
	JoinTypeManyToMany = "ManyToMany"
)

type Field struct {
	Name           string   `yaml:"name"`
	Alias          string   `yaml:"alias"`
	Description    string   `yaml:"description"`
	Tags           string   `yaml:"tags"`
	IsSkip         bool     `yaml:"isSkip"`
	IsRequired     bool     `yaml:"isRequired"`
	IsPrimaryKey   bool     `yaml:"isPrimaryKey"`
	IsSortable     bool     `yaml:"isSortable"`
	IsFilterable   bool     `yaml:"isFilterable"`
	JoinType       JoinType `yaml:"joinType"`
	TableName      string   `yaml:"tableName"`
	JoinTableName  string   `yaml:"JoinTableName"`
	ForeignKey     string   `yaml:"foreignKey"`
	References     string   `yaml:"references"`
	JoinForeignKey string   `yaml:"joinForeignKey"`
	JoinReferences string   `yaml:"joinReferences"`
	Operations     []string `yaml:"operations"`

	MetaType  MetaType
	TagsHTML  template.HTML
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
	f.NameSnakePlural = plural(f.NameSnake)
	f.NameKebab = snake(name)
	f.NameKebabPlural = plural(f.NameKebab)
	f.NameCamel = pascal(name)
	f.NameCamelPlural = plural(f.NameCamel)
	f.NameLowerCamel = camel(name)
	f.NameLowerCamelPlural = plural(f.NameLowerCamel)
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

	// fmt.Printf("columnType: %-20s %-20s %d\n", columnType, columnType, maxLength)
	return maxLength
}

func shiftMetaField(column *Column, metaTypes map[string]MetaType) *Field {

	isUnsigned := strings.Contains(column.Type, " unsigned") || strings.Contains(column.Type, " UNSIGNED")

	columnType := strings.ToLower(column.DataType)
	if isUnsigned {
		columnType = "u" + columnType
	}

	maxLength := parseColumnLength(column.Type)

	field := &Field{
		Name:        column.Name,
		Description: column.Comment,
		MetaType:    metaTypes[columnType], //meta
		MaxLength:   maxLength,
	}
	// fmt.Printf("shiftMetaField: %s -- %s -- %s -- %s -- %s -- %s -- %#v\n", column.Name, column.Type, column.DataType, column.Comment, columnType, metaTypes[columnType], field)

	field.genName()

	fmt.Println("shiftMetaField:", column, field.MetaType, isUnsigned, columnType)
	if strings.ToUpper(column.Key) == "PRI" {
		field.IsPrimaryKey = true
		// field.MetaType.GqlType = "ID"
		field.IsRequired = true
	}

	tags := `json:"` + field.NameLowerCamel + `" gorm:"` + field.NameLowerCamel
	if field.MaxLength > 0 {
		tags += ";size:" + strconv.Itoa(field.MaxLength)
	}
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

	return field
}

func mergeField(field *Field, fieldInCfg *Field) *Field {

	if fieldInCfg == nil {
		return field
	}

	if fieldInCfg.Alias != "" {
		field.Alias = fieldInCfg.Alias
	}

	if fieldInCfg.Description != "" {
		field.Description = fieldInCfg.Description
	}
	if fieldInCfg.Tags != "" {
		field.Tags = fieldInCfg.Tags
	}
	if fieldInCfg.JoinType != "" {
		field.JoinType = fieldInCfg.JoinType
	}
	if fieldInCfg.IsSkip {
		field.IsSkip = fieldInCfg.IsSkip
	}
	if fieldInCfg.IsRequired {
		field.IsRequired = fieldInCfg.IsRequired
	}
	if fieldInCfg.IsPrimaryKey {
		field.IsPrimaryKey = fieldInCfg.IsPrimaryKey
	}
	if fieldInCfg.IsSortable {
		field.IsSortable = fieldInCfg.IsSortable
	}
	if fieldInCfg.IsFilterable {
		field.IsFilterable = fieldInCfg.IsFilterable
		field.Operations = fieldInCfg.Operations
	}
	if fieldInCfg.ForeignKey != "" {
		field.ForeignKey = fieldInCfg.ForeignKey
	}
	if fieldInCfg.References != "" {
		field.References = fieldInCfg.References
	}
	if fieldInCfg.JoinForeignKey != "" {
		field.JoinForeignKey = fieldInCfg.JoinForeignKey
	}
	if fieldInCfg.JoinReferences != "" {
		field.JoinReferences = fieldInCfg.JoinReferences
	}

	return field
}
