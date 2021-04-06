package internal

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{ .Extra.pkgpath }}/pkg/app"
	"{{ .Extra.pkgpath }}/pkg/transports/http"
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
