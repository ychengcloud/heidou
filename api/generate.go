package api

import (
	"github.com/decker502/heidou"
	"github.com/decker502/heidou/loader"
)

func Generate(cfg *heidou.Config) error {
	o := &loader.Options{
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		Charset:  cfg.DBConfig.Charset,
	}
	loader := loader.NewMysqlSchemaLoader(o, cfg.DBConfig.Name)
	g := heidou.NewGenerator(cfg, loader)
	return g.Generate()
}
