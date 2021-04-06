package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{ .Extra.pkgpath }}/pkg/transports/http"
)

type Application struct {
	name       string
	logger     *zap.Logger
	httpServer *http.Server
}

type Option func(app *Application) error

func HttpServerOption(svr *http.Server) Option {
	return func(app *Application) error {
		svr.Application(app.name)
		app.httpServer = svr

		return nil
	}
}

func New(name string, logger *zap.Logger, options ...Option) (*Application, error) {
	app := &Application{
		name:   name,
		logger: logger.With(zap.String("type", "Application")),
	}

	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *Application) Start() error {
	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return errors.Wrap(err, "http server start error")
		}
	}

	return nil
}

func (a *Application) AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	s := <-c
	a.logger.Info("receive a signal", zap.String("signal", s.String()))
	if a.httpServer != nil {
		if err := a.httpServer.Stop(); err != nil {
			a.logger.Warn("stop http server error", zap.Error(err))
		}
	}

	os.Exit(0)

}
