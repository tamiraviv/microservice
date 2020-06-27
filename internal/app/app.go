package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"microservice/internal/pkg/errors"

	log "github.com/sirupsen/logrus"
)

const (
	confKeyLogBase  = "log"
	confKeyLogLevel = confKeyLogBase + ".level"
)

// Configuration expose an interface of configuration related actions
type Configuration interface {
	Get(key string) interface{}
	GetString(key string) (string, error)
	GetBool(key string) (bool, error)
	GetInt(key string) (int, error)
	GetDuration(key string) (time.Duration, error)
	IsSet(key string) bool
}

// RestServer defines a driving adapter interface
type RestServer interface {
	Start() error
	Stop(context.Context) error
}

// App defines the application struct
type App struct {
	restServer RestServer
}

// NewApp returns a new instance of the App struct
func NewApp(conf Configuration, rs RestServer) (*App, error) {
	logLevel, err := conf.GetString(confKeyLogLevel)
	if err != nil {
		return nil, err
	}

	logrusLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse log level (%s)", logLevel)
	}

	log.SetLevel(logrusLevel)

	return &App{
		restServer: rs,
	}, nil
}

// Start begins the flow of the app
func (a *App) Start() error {
	sig := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := a.restServer.Start(); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return errors.Wrap(err, "Failed to start drivers")
	case s := <-sig:
		log.Infof("Got (%s) signal to terminate application", s.String())
		return nil
	}
}

// Stop does a graceful shutdown of the app
func (a *App) Stop(ctx context.Context) error {
	if err := a.restServer.Stop(ctx); err != nil {
		return errors.Wrap(err, "Failed to gracefully stop rest server")
	}

	return nil
}
