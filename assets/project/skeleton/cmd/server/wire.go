// +build wireinject

package main

import (
	"github.com/google/wire"

	"{{ . }}/internal"
	"{{ . }}/internal/controllers"
	genControllers "{{ . }}/internal/gen/controllers"
	genRepositories "{{ . }}/internal/gen/repositories"
	genServices "{{ . }}/internal/gen/services"
	"{{ . }}/internal/repositories"
	"{{ . }}/internal/services"
	"{{ . }}/pkg/app"
	"{{ . }}/pkg/auth"
	"{{ . }}/pkg/auth/casbin"
	"{{ . }}/pkg/config"
	"{{ . }}/pkg/database"
	"{{ . }}/pkg/jaeger"
	"{{ . }}/pkg/log"
	"{{ . }}/pkg/transports/http"
	"{{ . }}/pkg/transports/http/middlewares/jwt"
	"{{ . }}/pkg/transports/http/middlewares/permission"
	"{{ . }}/pkg/validator"
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
	http.New, http.NewRouter,
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
