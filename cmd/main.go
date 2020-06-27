package main

import (
	"context"
	"os"
	"time"

	"microservice/wire"

	log "github.com/sirupsen/logrus"
)

const (
	initializationTimeout = 15
	stopTimeout           = 10
	exitCauseDone         = 0
	exitCauseError        = 1
)

func main() {
	initCtx, initCtxCancel := context.WithTimeout(context.Background(), initializationTimeout*time.Second)

	app, err := wire.InitializeApplication(initCtx)
	if err != nil {
		log.Errorf("Failed to inject dependencies: %s", err)
		initCtxCancel()
		os.Exit(exitCauseError)
	}

	exitCode := exitCauseDone
	if err := app.Start(); err != nil {
		log.Errorf("Failed to start application: %s", err)
		exitCode = exitCauseError
	}

	stopCtx, stopCtxCancel := context.WithTimeout(context.Background(), stopTimeout*time.Second)
	if err := app.Stop(stopCtx); err != nil {
		log.Errorf("Failed to gracefully stop application: %s", err)
		exitCode = exitCauseError
	}

	log.Info("Application closed")
	stopCtxCancel()
	os.Exit(exitCode)
}
