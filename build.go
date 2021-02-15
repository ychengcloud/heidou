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

	"golang.org/x/tools/imports"
)

type templateNode struct {
	NameFormat string
	FileName   string
}

var goModBase = templateNode{
	NameFormat: "go.mod",
	FileName:   "templates/go.mod.tmpl",
}

var modelsBase = templateNode{
	NameFormat: "internal/gen/models/models.go",
	FileName:   "templates/models_base.tmpl",
}

var controllersBase = templateNode{
	NameFormat: "internal/gen/controllers/controllers.go",
	FileName:   "templates/controllers_base.tmpl",
}

var repositoriesBase = templateNode{
	NameFormat: "internal/gen/repositories/repositories.go",
	FileName:   "templates/repositories_base.tmpl",
}

var servicesBase = templateNode{
	NameFormat: "internal/gen/services/services.go",
	FileName:   "templates/services_base.tmpl",
}

var controllers = templateNode{
	NameFormat: "internal/gen/controllers/%s.go",
	FileName:   "templates/controllers.tmpl",
}

var models = templateNode{
	NameFormat: "internal/gen/models/%s.go",
	FileName:   "templates/models.tmpl",
}

var repositories = templateNode{
	NameFormat: "internal/gen/repositories/%s.go",
	FileName:   "templates/repositories.tmpl",
}

var services = templateNode{
	NameFormat: "internal/gen/services/%s.go",
	FileName:   "templates/services.tmpl",
}

var swagger = templateNode{
	NameFormat: "assets/doc/swagger.yaml",
	FileName:   "templates/swagger.tmpl",
}

var assetsGenerate = templateNode{
	NameFormat: "cmd/server/assets_generate.go",
	FileName:   "templates/assets_generate.tmpl",
}

var parseBaseList = []templateNode{
	goModBase,
	modelsBase,
	controllersBase,
	repositoriesBase,
	servicesBase,
	swagger,
	assetsGenerate,
}

var parseRepeatList = []templateNode{
	controllers,
	models,
	repositories,
	services,
}

// format
func format(filename string, content []byte) error {
	ext := filepath.Ext(filename)
	data := content
	if ext == ".go" {
		var err error
		data, err = imports.Process(filename, content, nil)
		if err != nil {
			return fmt.Errorf("format file %s: %v", filename, err)
		}

	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("write file %s: %v", filename, err)
	}
	return nil
}

func (n *templateNode) ParseExecute(dir fs.FS, pathArg string, data interface{}) error {
	var path string
	if pathArg != "" {
		path = fmt.Sprintf(n.NameFormat, pathArg)
	} else {
		path = n.NameFormat
	}

	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	err := os.MkdirAll(filepath.Dir(path), 0744)
	if err != nil {
		return err
	}

	name := filepath.Base(n.FileName)
	t := template.New(name).Funcs(Funcs)
	tmpl, err := t.ParseFS(dir, n.FileName)
	if err != nil {
		return err
	}

	b := bytes.NewBuffer(nil)

	if err := tmpl.Execute(b, data); err != nil {
		return err
	}

	if err := format(path, b.Bytes()); err != nil {
		return err
	}
	return nil
}

// suffix不为空时，去掉生成文件的匹配后缀名
func build(dir fs.FS, root, dest string, trimSuffix bool, data interface{}) error {
	
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
				return err
			}
		} else {
			t := template.New(filepath.Base(path)).Funcs(Funcs)
			tmpl, err := t.ParseFS(dir, path)
			if err != nil {
				return err
			}

			if trimSuffix {
				suffix := filepath.Ext(target)
				target = strings.TrimSuffix(target, suffix)
			}

			b := bytes.NewBuffer(nil)
			if err := tmpl.Execute(b, data); err != nil {
				return err
			}

			if err := format(target, b.Bytes()); err != nil {
				return err
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
