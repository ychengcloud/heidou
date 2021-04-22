package api

import (
	"strings"

	"github.com/ychengcloud/heidou"
	"github.com/ychengcloud/heidou/loader"
)

func Generate(cfg *heidou.Config) error {
	var loaderInstance heidou.Loader
	switch strings.ToLower(strings.Trim(cfg.Loader, " ")) {
	case "mysql":
		o := &loader.MysqlSchemaLoaderOptions{
			User:     cfg.DBConfig.User,
			Password: cfg.DBConfig.Password,
			Host:     cfg.DBConfig.Host,
			Port:     cfg.DBConfig.Port,
			Charset:  cfg.DBConfig.Charset,
		}
		loaderInstance = loader.NewMysqlSchemaLoader(o, cfg.DBConfig.Name)
	default:
		loaderInstance = loader.NewMemoryLoader()
	}

	g, err := heidou.NewGenerator(cfg, loaderInstance)
	if err != nil {
		return err
	}
	return g.Generate()
}
