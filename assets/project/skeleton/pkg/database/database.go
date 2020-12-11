package database

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Options is  configuration of database
type Options struct {
	Dialect  string `yaml:"dialect"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Charset  string `yaml:"charset"`
	Debug    bool
}

// Init 初始化数据库
func New(o *Options) (*gorm.DB, error) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", o.User, o.Password, o.Host, o.Port, o.Name, o.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "gorm open database connection error")
	}

	if o.Debug {
		db = db.Debug()
	}

	return db, nil
}
