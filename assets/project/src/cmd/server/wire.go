// +build wireinject

package main

import (
	"github.com/google/wire"

	server "{{ . }}/internal"
	"{{ . }}/internal/handlers"
	"{{ . }}/pkg/app"
	"{{ . }}/pkg/config"
	"{{ . }}/pkg/database"
	"{{ . }}/pkg/log"
	"{{ . }}/pkg/transports/http"
	"{{ . }}/pkg/validator"
)

var providerSet = wire.NewSet(
	NewHttpOptions,
	NewServerOptions,
	NewHandlersOptions,
	NewLogOptions,
	NewDatabaseOptions,
	log.New,
	config.New,
	// jaeger.New, jaeger.NewConfiguration,
	database.New,
	validator.New,
	http.New, http.NewRouter,
	handlers.CreateGqlHandlers,
	server.NewServer,
)

func CreateApp(cf string) (*app.Application, error) {
	panic(wire.Build(providerSet))
}
