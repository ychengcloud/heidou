package heidou

import (
	"fmt"
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
	Name            string   `yaml:"name"`
	Alias           string   `yaml:"alias"`
	Description     string   `yaml:"description"`
	Tags            string   `yaml:"tags"`
	IsSkip          bool     `yaml:"isSkip"`
	IsRequired      bool     `yaml:"isRequired"`
	IsPrimaryKey    bool     `yaml:"isPrimaryKey"`
	IsForeignKey    bool     `yaml:"isForeignKey"`
	IsAutoIncrement bool     `yaml:"isAutoIncrement"`
	IsSortable      bool     `yaml:"isSortable"`
	IsFilterable    bool     `yaml:"isFilterable"`
	JoinType        JoinType `yaml:"joinType"`
	TableName       string   `yaml:"tableName"`
	JoinTableName   string   `yaml:"JoinTableName"`
	ForeignKey      string   `yaml:"foreignKey"`
	References      string   `yaml:"references"`
	JoinForeignKey  string   `yaml:"joinForeignKey"`
	JoinReferences  string   `yaml:"joinReferences"`
	Operations      []string `yaml:"operations"`

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
	}

	field.genName()
	field.handleTags()

	return field
}

func (f *Field) handleTags() {
	tags := `json:"` + f.NameLowerCamel + `" gorm:"column:` + f.NameSnake
	if f.IsPrimaryKey {
		tags += ";primaryKey"
	}
	if f.IsAutoIncrement {
		tags += ";autoIncrement"
	}
	// 设置了 default tag, gorm的createdat updatedat 逻辑会失效
	// if f.MetaType.GoType == "time.Time" {
	// 	tags += ";default: '1970-01-01 00:00:00'"
	// }
	if f.MaxLength > 0 {
		tags += ";size:" + strconv.Itoa(f.MaxLength)
	}
	tags += `"`
	if len(f.Tags) > 0 {
		tags = tags + " " + f.Tags
	}

	// fmt.Println("tags:", tags)
	f.TagsHTML = template.HTML(tags)
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
	if fieldInCfg.JoinType == "" {
		fieldInCfg.JoinType = JoinTypeNone
	}

	field.JoinType = fieldInCfg.JoinType

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

	field.handleTags()

	return field
}

func (f *Field) HandleAssociation() {
	if f.JoinType == JoinTypeNone {
		return
	}

	tags := `json:"` + f.NameLowerCamel + `" gorm:"` + f.NameLowerCamel

	if f.JoinType == JoinTypeManyToMany || f.JoinType == JoinTypeHasMany {
		tags = `json:"` + f.NameLowerCamelPlural + `" gorm:"` + f.NameLowerCamelPlural

	}
	if f.JoinTableName != "" {
		tags += ";many2many:" + f.JoinTableName
	}
	if f.ForeignKey != "" {
		tags += ";foreignKey:" + f.ForeignKey
	}
	if f.References != "" {
		tags += ";references:" + f.References
	}
	if f.JoinForeignKey != "" {
		tags += ";joinForeignKey:" + f.JoinForeignKey
	}
	if f.JoinReferences != "" {
		tags += ";joinReferences:" + f.JoinReferences
	}
	tags += `"`

	fmt.Println("tags:", tags)
	f.TagsHTML = template.HTML(tags)
}
