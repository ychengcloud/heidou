package heidou

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/shurcooL/httpfs/text/vfstemplate"
	"github.com/shurcooL/httpfs/vfsutil"

	"github.com/decker502/heidou/assets"
)

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

// suffix不为空时，去掉生成文件的匹配后缀名
func build(fs http.FileSystem, root, dest string, trimSuffix bool, data interface{}) error {
	walkFn := func(path string, fi os.FileInfo, err error) error {
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

		if fi.IsDir() {
			err = os.MkdirAll(target, 0744)
			if err != nil {
				return err
			}
		} else {
			t := template.New(filepath.Base(path)).Funcs(Funcs)
			tmpl, err := vfstemplate.ParseFiles(fs, t, path)
			if err != nil {
				return err
			}

			if trimSuffix {
				suffix := filepath.Ext(target)
				target = strings.TrimSuffix(target, suffix)
			}
			file, err := os.Create(target)
			if err != nil {
				return err
			}
			defer file.Close()

			err = tmpl.Execute(file, data)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := vfsutil.Walk(fs, root, walkFn)
	if err != nil {
		return err
	}

	return nil
}
