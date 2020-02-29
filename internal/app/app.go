package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"microservice/internal/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type RestServer interface {
	Start() error
	Stop(context.Context) error
}

type App struct {
	restServer RestServer
}

func NewApp(rs RestServer) App {
	return App{
		restServer: rs,
	}
}

func (a *App) Start() error {
	sig := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errChan <- a.restServer.Start()
	}()

	select {
	case err := <-errChan:
		return errors.Wrap(err, "Failed to start drivers")
	case s := <-sig:
		log.Infof("Got (%s) signal to terminate application", s.String())
		return nil
	}
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.restServer.Stop(ctx); err != nil {
		return errors.Wrap(err, "Failed to gracefully stop rest server")
	}

	return nil
}
