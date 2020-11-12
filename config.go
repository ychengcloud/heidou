package heidou

import "text/template"

type Config struct {
	PkgPath  string   `yaml:"pkgPath"`
	DBConfig DBConfig `yaml:"db" mapstructure:"db"`
	Tables   []*Table `yaml:"tables"`

	Funcs template.FuncMap
}

type DBConfig struct {
	// 表名前缀
	Prefix   string `yaml:"prefix"`
	Dialect  string `yaml:"dialect"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Name     string `yaml:"name"`
	Charset  string `yaml:"charset"`
}
