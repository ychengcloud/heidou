package heidou

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/horcus/heidou/assets"
	"github.com/shurcooL/httpfs/vfsutil"
)

type Generator struct {
	Config    *Config
	Data      *Data
	Loader    Loader
	MetaTypes map[string]MetaType
}

type Data struct {
	PkgPath         string
	HasTimeField    bool
	IsImportStrings bool

	Tables []*Table
}

func NewGenerator(cfg *Config, loader Loader) *Generator {
	data := &Data{
		PkgPath: cfg.PkgPath,
	}

	gen := &Generator{
		Config: cfg,
		Data:   data,
		Loader: loader,
	}

	mappings := must(loadMappings("/mappings.yaml"))

	gen.MetaTypes = make(map[string]MetaType, len(mappings.MetaTypes))
	for _, v := range mappings.MetaTypes {
		gen.MetaTypes[v.SQLType] = v
	}
	return gen
}

// getTable 根据表名，查找表对象
func (g *Generator) getTableInCfg(name string) *Table {

	for _, table := range g.Config.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

// getTable 根据表名，查找表对象
func (g *Generator) getTable(name string) *Table {

	for _, table := range g.Data.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

//生成many2many的反向引用
func (g *Generator) handleBackReference(table *Table, backReferenceTable *Table) {

	backReferenceInfo := &BackReferenceInfo{
		Name:                 table.Name,
		NameCamel:            table.NameCamel,
		NameCamelPlural:      table.NameCamelPlural,
		NameLowerCamelPlural: table.NameLowerCamelPlural,
		JoinTableName:        table.Name,
	}

	//配置了此表的反向多对多关联
	find := false
	for _, fieldInCfg := range backReferenceTable.Fields {
		if fieldInCfg.JoinType == JoinTypeManyToMany && fieldInCfg.TableName == table.Name {
			find = true
			break
		}
	}
	if !find {
		fmt.Println("BackReferenceInfos:", backReferenceTable.Name, backReferenceInfo)
		backReferenceTable.BackReferenceInfos = append(backReferenceTable.BackReferenceInfos, backReferenceInfo)
	}

	return
}

// 所有元表信息生成后，再处理关联字段的模型信息
func (g *Generator) handleAssociation() error {
	for _, table := range g.Data.Tables {
		for _, field := range table.Fields {
			if field.JoinType == JoinTypeNone {
				continue
			}
			field.JoinTable = g.getTable(field.TableName)
			if field.JoinTable == nil {
				return fmt.Errorf("Something wrong, can't find %s", field)
			}
			if field.JoinType == JoinTypeManyToMany {
				g.handleBackReference(table, field.JoinTable)
			}

		}
	}
	return nil
}

// 根据配置信息和数据库表结构生成 Graphql 相关代码
func (g *Generator) Generate() error {
	metaTables, err := g.Loader.LoadMetaTable()
	if err != nil {
		return err
	}

	for _, metaTable := range metaTables {
		// if metaTable.Name != "product" && metaTable.Name != "category" && metaTable.Name != "product_variant" {
		// 	continue
		// }
		tableInCfg := g.getTableInCfg(metaTable.Name)
		table := MergeTable(metaTable, tableInCfg, g.MetaTypes)
		if table == nil || table.IsSkip {
			continue
		}

		table.genName()
		table.PkgPath = g.Data.PkgPath
		if table.HasTimeField {
			g.Data.HasTimeField = true
		}
		if table.IsImportStrings {
			g.Data.IsImportStrings = true
		}
		g.Data.Tables = append(g.Data.Tables, table)
	}

	if err := g.handleAssociation(); err != nil {
		return err
	}
	if err := g.gen(); err != nil {
		return err
	}
	return nil
}

func GenProject(dest string, pkgPath string) error {

	err := genSkeleton(dest, pkgPath)
	if err != nil {
		return err
	}

	return nil
}

func genSkeleton(dest string, data interface{}) error {

	err := build(assets.Project, "/skeleton", dest, false, data)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) gen() error {
	for _, node := range parseBaseList {
		err := node.ParseExecute(assets.Project, "", g.Data)
		if err != nil {
			return fmt.Errorf("parse [%s] template failed with error : %s", node.NameFormat, err)
		}
	}

	for _, table := range g.Data.Tables {
		tableName := table.Name
		//generate model from table
		for _, node := range parseRepeatList {
			err := node.ParseExecute(assets.Project, tableName, table)
			if err != nil {
				return fmt.Errorf("parse [%s] template failed with error : %s", node.NameFormat, err)
			}
		}
	}
	return nil
}
func must(sm *MetaTypes, err error) *MetaTypes {
	if err != nil {
		panic(err)
	}
	return sm
}

// 加载 SQL 类型映射文件数据
func loadMappings(mappingFileName string) (*MetaTypes, error) {
	byteValue, err := vfsutil.ReadFile(assets.Project, mappingFileName)
	if err != nil {
		fmt.Printf("Error loading mapping file %s error: %v\n", mappingFileName, err)
		return nil, err
	}

	return processMappings(byteValue)
}

// 处理 SQL 类型映射
func processMappings(mappingContent []byte) (*MetaTypes, error) {
	mappings := &MetaTypes{}
	err := yaml.Unmarshal(mappingContent, mappings)
	if err != nil {
		fmt.Printf("Error unmarshalling yaml error: %v\n", err)
		return nil, err
	}

	// fmt.Println("mappings:", string(mappingContent), mappings)
	// fmt.Println("mappings:", mappings)
	return mappings, nil
}
