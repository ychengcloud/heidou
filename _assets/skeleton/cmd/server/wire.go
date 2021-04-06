// +build wireinject

package main

import (
	"github.com/google/wire"

	"{{ .Extra.pkgpath }}/internal"
	"{{ .Extra.pkgpath }}/internal/controllers"
	genControllers "{{ .Extra.pkgpath }}/internal/gen/controllers"
	genRepositories "{{ .Extra.pkgpath }}/internal/gen/repositories"
	genServices "{{ .Extra.pkgpath }}/internal/gen/services"
	"{{ .Extra.pkgpath }}/internal/repositories"
	"{{ .Extra.pkgpath }}/internal/services"
	"{{ .Extra.pkgpath }}/pkg/app"
	"{{ .Extra.pkgpath }}/pkg/auth"
	"{{ .Extra.pkgpath }}/pkg/auth/casbin"
	"{{ .Extra.pkgpath }}/pkg/config"
	"{{ .Extra.pkgpath }}/pkg/database"
	"{{ .Extra.pkgpath }}/pkg/jaeger"
	"{{ .Extra.pkgpath }}/pkg/log"
	"{{ .Extra.pkgpath }}/pkg/transports/http"
	"{{ .Extra.pkgpath }}/pkg/transports/http/middlewares/jwt"
	"{{ .Extra.pkgpath }}/pkg/transports/http/middlewares/permission"
	"{{ .Extra.pkgpath }}/pkg/validator"
)

var providerSet = wire.NewSet(
	config.New,
	NewHttpOptions, http.New, http.NewRouter,
	NewServerOptions, internal.NewServer,
	NewLogOptions, log.New,
	NewDatabaseOptions, database.New,
	NewAuthOptions, auth.ProviderSet,
	NewJaegerConfig, jaeger.New,
	validator.New,
	jwt.ProviderSet,
	permission.ProviderSet,
	casbin.ProviderSet,
	genRepositories.BaseProviderSet,
	genControllers.BaseProviderSet,
	genServices.BaseProviderSet,
	repositories.ProviderSet,
	controllers.ProviderSet,
	services.ProviderSet,
	internal.NewJWTCallback,
	internal.NewAuthCallback,
)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
