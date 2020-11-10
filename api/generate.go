package api

import (
	"github.com/decker502/heidou/internal/gen"
	"github.com/decker502/heidou/internal/loader"
)

func Generate(cfg *gen.Config) error {
	o := &loader.Options{
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		Charset:  cfg.DBConfig.Charset,
	}
	loader := loader.NewMysqlSchemaLoader(o, cfg.DBConfig.Name)
	g := gen.NewGenerator(cfg, loader)
	return g.Generate()
}
