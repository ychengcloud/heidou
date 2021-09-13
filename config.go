package heidou

import "text/template"

type Config struct {
	ProjectName string      `yaml:"projectName"`
	Overwrite   bool        `yaml:"overwrite"`
	Loader      string      `yaml:"loader"`
	Extra       interface{} `yaml:"extra"`
	DBConfig    DBConfig    `yaml:"db" mapstructure:"db"`
	Tables      []*Table    `yaml:"tables"`
	Templates   *Template   `yaml:"templates" mapstructure:"templates"`
	Delim       Delim       `yaml:"delim"`

	TemplatesPath string
	Funcs         template.FuncMap
}

type Delim struct {
	Left  string `yaml:"left"`
	Right string `yaml:"right"`
}

// 生成的模板
type TempInfo struct {
	// NamePattern 目标路径，支持简单的字符串替换
	NamePattern string `yaml:"namePattern"`
	// Path 模板路径名，以 templates为相对路径
	Path string `yaml:"path"`

	//Type 模板类型，默认: table schema类型： 全部表数据生成一个文件  table 类型：一个数据表生成一个文件
	Type string `yaml:"type" mapstructure:"type"`
}
type Template struct {
	NameStyle string      `yaml:"nameStyle"`
	InfoList  []*TempInfo `yaml:"info" mapstructure:"info"`

	//引用的模板，只会被主模板引用，不会单独生成，支持golang glob pattern
	References []string `yaml:"references" mapstructure:"references"`
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
