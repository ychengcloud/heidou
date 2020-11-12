package heidou

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/shurcooL/httpfs/filter"
	"github.com/shurcooL/httpfs/text/vfstemplate"
	"github.com/shurcooL/httpfs/vfsutil"
	"honnef.co/go/tools/config"

	"github.com/decker502/heidou/assets"
)

const (
	GraphqlWeb = "graphql"
)

func GenProject(dest string, pkgPath string) error {
	err := genGraphqlWeb(dest)
	if err != nil {
		return err
	}

	err = genSrc(dest, pkgPath)
	if err != nil {
		return err
	}

	return nil
}

func GenGql(cfg *config.Config) error {
	// return genGql("/", data)
	return nil
}

func genGraphqlWeb(dest string) error {
	walkFn := func(path string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			return err
		}

		mask := syscall.Umask(0)
		defer syscall.Umask(mask)

		relPath, err := filepath.Rel("/src", path)
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
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}
			if err = ioutil.WriteFile(target, b, 0744); err != nil {
				return err
			}
		}

		return nil
	}

	fs := filter.Keep(assets.Project, func(path string, fi os.FileInfo) bool {
		// logrus.Println(" Keep path:", path)
		return path == "/" ||
			path == "/src" ||
			path == "/src/"+GraphqlWeb ||
			strings.HasPrefix(path, "/src/"+GraphqlWeb+"/")
		// return true
	})

	err := vfsutil.WalkFiles(fs, "/src", walkFn)
	if err != nil {
		return err
	}

	return nil
}

func genSrc(dest string, data interface{}) error {
	fs := filter.Skip(assets.Project, func(path string, fi os.FileInfo) bool {
		return path == "/src/"+GraphqlWeb ||
			strings.HasPrefix(path, "/src/"+GraphqlWeb+"/")
	})

	err := build(fs, "/src", dest, false, data)
	if err != nil {
		return err
	}
	return nil
}

func genGql(dest string, data interface{}) error {
	err := build(assets.Project, "/templates", dest, true, data)
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
			fmt.Println("test")

		}

		return nil
	}

	err := vfsutil.Walk(fs, root, walkFn)
	if err != nil {
		return err
	}

	return nil
}
