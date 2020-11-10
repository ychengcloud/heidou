package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	server "{{ . }}/internal"
	"{{ . }}/internal/handlers"
	"{{ . }}/pkg/database"
	"{{ . }}/pkg/log"
	"{{ . }}/pkg/transports/http"
)

func NewHttpOptions(v *viper.Viper) (*http.Options, error) {
	var (
		err error
		o   = new(http.Options)
	)

	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, err
	}

	return o, err
}

func NewServerOptions(v *viper.Viper, logger *zap.Logger) (*server.Options, error) {
	var err error
	o := new(server.Options)
	if err = v.UnmarshalKey("app", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal app option error")
	}

	logger.Info("load application options success")

	return o, err
}

func NewHandlersOptions(v *viper.Viper, logger *zap.Logger) (*handlers.Options, error) {
	var err error
	o := new(handlers.Options)
	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal app option error")
	}

	logger.Info("load application options success")

	return o, err
}

func NewLogOptions(v *viper.Viper) (*log.Options, error) {
	var (
		err error
		o   = new(log.Options)
	)

	if err = v.UnmarshalKey("log", o); err != nil {
		return nil, err
	}
	fmt.Println("new log:", o.Filename)
	return o, err
}

func NewDatabaseOptions(v *viper.Viper, logger *zap.Logger) (*database.Options, error) {
	var err error
	o := new(database.Options)
	if err = v.UnmarshalKey("db", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}

	logger.Info("load database options success",
		zap.String("User", o.User),
		zap.String("Host", o.Host),
		zap.Int("Port", o.Port),
	)
	return o, err
}
