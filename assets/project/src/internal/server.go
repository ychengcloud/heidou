package server

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{ . }}/pkg/app"
	"{{ . }}/pkg/transports/http"
)

type Options struct {
	Name string
}

func NewServer(o *Options, logger *zap.Logger, hs *http.Server) (*app.Application, error) {
	a, err := app.New(o.Name, logger, app.HttpServerOption(hs))

	if err != nil {
		return nil, errors.Wrap(err, "new app error")
	}

	return a, nil
}
