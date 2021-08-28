package heidou

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

const (
	AssetsRoot   = "_assets"
	TmplBasePath = "templates/"
)

type Generator struct {
	Config    *Config
	Data      *Data
	Loader    Loader
	Assets    fs.FS
	MetaTypes map[string]MetaType
}

type Data struct {
	ProjectName     string
	Extra           interface{}
	HasTimeField    bool
	IsImportStrings bool

	Tables []*Table
}

func NewGenerator(cfg *Config, loader Loader) (*Generator, error) {
	data := &Data{
		ProjectName: cfg.ProjectName,
		Extra:       cfg.Extra,
	}

	gen := &Generator{
		Config: cfg,
		Data:   data,
		Loader: loader,
	}

	var err error

	if len(cfg.TemplatesPath) > 0 {
		gen.Assets = os.DirFS(cfg.TemplatesPath)
	} else {
		gen.Assets, err = fs.Sub(Assets, AssetsRoot)
		if err != nil {
			return nil, err
		}

	}

	mappings := must(loadMappings("mappings.yaml"))

	var metaTypes []MetaType
	for _, m := range mappings.MetaTypeInfos {
		if m.Dialect == cfg.Loader {
			metaTypes = m.MetaTypes
			break
		}
	}

	if len(metaTypes) == 0 {
		return nil, fmt.Errorf("%s Meta type mapping can't be found ", cfg.Loader)
	}

	gen.MetaTypes = make(map[string]MetaType, len(metaTypes))
	for _, v := range metaTypes {
		gen.MetaTypes[v.SQLType] = v
	}
	return gen, nil
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
func handleBackReference(table *Table, backReferenceTable *Table) {

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
		if fieldInCfg.JoinType == JoinTypeManyToMany && fieldInCfg.RefTableName == table.Name {
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

//关联外键
func handleAssociationForeignKey(table *Table, foreignKey string) {
	for _, field := range table.Fields {
		if field.NameCamel == foreignKey {
			field.IsForeignKey = true
		}
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
			table.HasJoinField = true
			field.RefTable = g.getTable(field.RefTableName)
			if field.RefTable == nil {
				return fmt.Errorf("Something wrong, can't find RefTable : %#v\n", field)
			}
			if field.JoinType == JoinTypeManyToMany {
				handleBackReference(table, field.RefTable)
			}

			foreignKey := table.NameCamel + "Id"
			if len(field.ForeignKey) > 0 {
				foreignKey = pascal(field.ForeignKey)
			}

			// 更新相应外键字段信息
			if field.JoinType == JoinTypeBelongTo {
				handleAssociationForeignKey(table, foreignKey)
			}

			if field.JoinType == JoinTypeHasOne || field.JoinType == JoinTypeHasMany {
				handleAssociationForeignKey(field.RefTable, foreignKey)
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

		//不支持联合主键或联合索引，如果表中定义了联合主键或联合索引，必须显式配置忽略相应的表
		if metaTable.hasCompositeKeys() {
			if tableInCfg == nil || !tableInCfg.IsSkip {
				return fmt.Errorf("Must be skip if it has composite keys: %s", metaTable.Name)
			}
		}
		table := MergeTable(metaTable, tableInCfg, g.MetaTypes)
		if table == nil || table.IsSkip {
			continue
		}

		table.genName()
		if table.HasTimeField {
			g.Data.HasTimeField = true
		}
		if table.IsImportStrings {
			g.Data.IsImportStrings = true
		}
		table.Extra = g.Data.Extra
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

// suffix不为空时，去掉生成文件的匹配后缀名
func (g *Generator) build(dir fs.FS, root, dest string, data interface{}, overwrite bool) error {

	walkFn := func(path string, entry fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		mask := Umask(0)
		defer Umask(mask)

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, relPath)
		// target := filepath.Join(dest, path)
		if entry.IsDir() {
			err = os.MkdirAll(target, 0744)
			if err != nil {
				return err
			}
		} else {
			var content []byte
			//如果后缀以 .tmpl .tpl结尾，则约定需要自动去除后缀
			suffix := filepath.Ext(path)
			if suffix == ".tmpl" || suffix == ".tpl" {
				target = strings.TrimSuffix(target, suffix)
				b := bytes.NewBuffer(nil)

				t := template.New(filepath.Base(path)).Funcs(Funcs).Funcs(sprig.GenericFuncMap())
				t.Delims(g.Config.Delim.Left, g.Config.Delim.Right)
				tmpl, err := t.ParseFS(dir, path)
				if err != nil {
					return err
				}

				if err := tmpl.Execute(b, data); err != nil {
					return err
				}

				content = b.Bytes()

			} else {
				content, err = fs.ReadFile(dir, path)
				if err != nil {
					return err
				}
			}

			// 根据 overwrite 标志决定是否覆盖已有文件
			if !overwrite {
				_, err := os.Stat(target)

				// 文件存在则返回，两种情况表示文件是存在的，1 err 为nil, 2 返回的err不包含不存在的标志
				if err == nil || !os.IsNotExist(err) {
					return err
				}
			}

			if err := os.WriteFile(target, content, 0644); err != nil {
				return err
			}

			suffix = filepath.Ext(target)
			if suffix == ".go" {

				if err := format(target, content); err != nil {
					return err
				}
			}
			return nil

		}
		return nil
	}

	// err := vfsutil.Walk(fs, root, walkFn)
	err := fs.WalkDir(dir, root, walkFn)
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) genSkeleton(assets fs.FS, dest string, data interface{}) error {
	err := g.build(assets, "skeleton", dest, data, g.Config.Overwrite)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) gen() error {
	err := g.genSkeleton(g.Assets, "./", g.Data)
	if err != nil {
		return err
	}

	for _, table := range g.Data.Tables {
		tableName := table.Name
		switch g.Config.TmplNameFormat {
		case "camel":
			tableName = table.NameCamel
		case "lowerCamel":
			tableName = table.NameLowerCamel
		case "camelPlural":
			tableName = table.NameCamelPlural
		case "lowerCamelPlural":
			tableName = table.NameLowerCamelPlural
		case "snake":
			tableName = table.NameSnake

		}

		//generate model from table
		for _, template := range g.Config.Templates {
			dest := fmt.Sprintf(template.NameFormat, tableName)
			tmplPath := template.Path

			//根据TypeName生成实际模板路径
			if table.TypeName != "" {
				paths := strings.SplitN(tmplPath, ".", 2)
				if paths != nil {
					tmplPath = paths[0] + "_" + table.TypeName
					if len(paths) == 2 {
						tmplPath = tmplPath + "." + paths[1]
					}
				} else {
					tmplPath = "_" + table.TypeName
				}

				tmplPath = TmplBasePath + tmplPath
				//不存在则使用默认的模板
				if _, err := fs.Stat(g.Assets, tmplPath); os.IsNotExist(err) {
					tmplPath = TmplBasePath + template.Path
				}
			} else {
				tmplPath = TmplBasePath + tmplPath
				//不存在则不生成
				if _, err := fs.Stat(g.Assets, tmplPath); os.IsNotExist(err) {
					continue
				}
			}

			err := g.build(g.Assets, tmplPath, dest, table, g.Config.Overwrite)
			if err != nil {
				return fmt.Errorf("parse [%s] template failed with error : %s", table.Name, err.Error())
			}

			// err := template.ParseExecute(g.Assets, tableName, table)
			// if err != nil {
			// 	return fmt.Errorf("parse [%s] template failed with error : %s", template.NameFormat, err)
			// }
		}
	}
	return nil
}
func must(sm *MetaTypeMapping, err error) *MetaTypeMapping {
	if err != nil {
		panic(err)
	}
	return sm
}

// 加载 SQL 类型映射文件数据
func loadMappings(mappingFileName string) (*MetaTypeMapping, error) {
	sub, err := fs.Sub(Assets, AssetsRoot)
	if err != nil {
		return nil, err
	}

	byteValue, err := fs.ReadFile(sub, mappingFileName)
	if err != nil {
		fmt.Printf("Error loading mapping file %s error: %v\n", mappingFileName, err)
		return nil, err
	}

	return processMappings(byteValue)
}

// 处理 SQL 类型映射
func processMappings(mappingContent []byte) (*MetaTypeMapping, error) {
	mappings := &MetaTypeMapping{}
	err := yaml.Unmarshal(mappingContent, mappings)
	if err != nil {
		fmt.Printf("Error unmarshalling yaml error: %v\n", err)
		return nil, err
	}

	// fmt.Println("mappings:", string(mappingContent), mappings)
	// fmt.Println("mappings:", mappings)
	return mappings, nil
}
