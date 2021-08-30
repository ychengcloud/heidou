package heidou

import "text/template"

type Config struct {
	ProjectName    string      `yaml:"projectName"`
	Overwrite      bool        `yaml:"overwrite"`
	Loader         string      `yaml:"loader"`
	Extra          interface{} `yaml:"extra"`
	DBConfig       DBConfig    `yaml:"db" mapstructure:"db"`
	Tables         []*Table    `yaml:"tables"`
	Templates      []*Template `yaml:"templates"`
	Delim          Delim       `yaml:"delim"`
	TmplNameFormat string      `yaml:"tmplNameFormat"`

	TemplatesPath string
	Funcs         template.FuncMap
}

type Delim struct {
	Left  string `yaml:"left"`
	Right string `yaml:"right"`
}

type Template struct {
	NameFormat string `yaml:"nameFormat"`
	Path       string `yaml:"path"`

	//Type 模板类型，默认: table schema类型： 全部表数据生成一个文件  table 类型：一个数据表生成一个文件
	Type string `yaml:"type" mapstructure:"type"`
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
