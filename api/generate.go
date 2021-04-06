package api

import (
	"github.com/horcus/heidou"
	"github.com/horcus/heidou/loader"
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
	g, err := heidou.NewGenerator(cfg, loader)
	if err != nil {
		return err
	}
	return g.Generate()
}
