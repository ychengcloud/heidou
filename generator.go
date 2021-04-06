package heidou

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"gopkg.in/yaml.v2"
)

const (
	AssetsRoot = "_assets"
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

	gen.MetaTypes = make(map[string]MetaType, len(mappings.MetaTypes))
	for _, v := range mappings.MetaTypes {
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
				return fmt.Errorf("Something wrong, can't find1 : %#v\n", field)
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
func (g *Generator) build(dir fs.FS, root, dest string, data interface{}) error {

	walkFn := func(path string, entry fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		mask := syscall.Umask(0)
		defer syscall.Umask(mask)

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, relPath)
		// target := filepath.Join(dest, path)
		if entry.IsDir() {
			err = os.MkdirAll(target, 0744)
			if err != nil {
				fmt.Println("parse fs1:", dir, path)
				return err
			}
		} else {
			t := template.New(filepath.Base(path)).Funcs(Funcs)
			t.Delims(g.Config.Delim.Left, g.Config.Delim.Right)
			tmpl, err := t.ParseFS(dir, path)
			if err != nil {
				fmt.Println("parse fs2:", dir, path)
				return err
			}

			//如果后缀以 .tmpl .tpl结尾，则约定需要自动去除后缀
			suffix := filepath.Ext(target)
			if suffix == ".tmpl" || suffix == ".tpl" {
				target = strings.TrimSuffix(target, suffix)
			}

			b := bytes.NewBuffer(nil)
			if err := tmpl.Execute(b, data); err != nil {
				fmt.Println("parse fs3:", dir, path)
				return err
			}

			if err := format(target, b.Bytes()); err != nil {
				fmt.Println("parse fs4:", dir, path)
				return err
			}
			return nil

		}
		return nil
	}

	// err := vfsutil.Walk(fs, root, walkFn)
	err := fs.WalkDir(dir, root, walkFn)
	if err != nil {
		fmt.Printf("parse fs5: %#v %s %#v\n", dir, root, err)

		return err
	}

	return nil
}

func (g *Generator) genSkeleton(assets fs.FS, dest string, data interface{}) error {
	err := g.build(assets, "skeleton", dest, data)
	if err != nil {
		fmt.Println("err:", err, dest)
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
		//generate model from table
		for _, template := range g.Config.Templates {
			path := fmt.Sprintf(template.NameFormat, tableName)

			err := g.build(g.Assets, "templates/"+template.Path, path, table)
			if err != nil {
				fmt.Println("err:", err, path)
				return err
			}

			// err := template.ParseExecute(g.Assets, tableName, table)
			// if err != nil {
			// 	return fmt.Errorf("parse [%s] template failed with error : %s", template.NameFormat, err)
			// }
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
