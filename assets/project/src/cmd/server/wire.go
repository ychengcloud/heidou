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
	"{{ . }}/pkg/log"
	"{{ . }}/pkg/transports/http"
	"{{ . }}/pkg/transports/http/middlewares/jwt"
	"{{ . }}/pkg/validator"
)

var providerSet = wire.NewSet(
	NewHttpOptions,
	NewServerOptions,
	NewLogOptions,
	NewDatabaseOptions,
	NewAuthOptions,
	log.New,
	config.New,
	// jaeger.New, jaeger.NewConfiguration,
	database.New,
	validator.New,
	auth.ProviderSet,
	jwt.ProviderSet,
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
	internal.NewServer,
)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
