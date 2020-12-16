package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"

	"{{ . }}/internal"
	"{{ . }}/pkg/auth"
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

func NewServerOptions(v *viper.Viper, logger *zap.Logger) (*internal.Options, error) {
	var err error
	o := new(internal.Options)
	if err = v.UnmarshalKey("app", o); err != nil {
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

func NewAuthOptions(v *viper.Viper, logger *zap.Logger) (*auth.Options, error) {
	var (
		err error
		o   = &auth.DefaultOptions
	)
	if err = v.UnmarshalKey("jwt", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal jwt options error")
	}

	o.Logger = logger
	logger.Info("load jwt options success",
		zap.String("ClaimsKey", o.ClaimsKey),
		zap.String("SigningKey", o.SigningKey),
	)

	var account string
	if err = v.UnmarshalKey("app.rootAccount", &account); err != nil {
		return nil, err
	}
	o.RootAccount = account
	
	return o, nil
}

func NewJaegerConfig(v *viper.Viper, logger *zap.Logger) (*config.Configuration, error) {
	var (
		err error
		c   = new(config.Configuration)
	)

	if err = v.UnmarshalKey("jaeger", c); err != nil {
		return nil, errors.Wrap(err, "unmarshal jaeger configuration error")
	}

	logger.Info("load jaeger configuration success")

	return c, nil
}
