package main

import (
	"context"
	"time"

	"microservice/wire"

	log "github.com/sirupsen/logrus"
)

const (
	initializationTimeout = 15
	stopTimeout           = 10
)

func main() {
	initCtx, initCtxCancel := context.WithTimeout(context.Background(), initializationTimeout*time.Second)
	defer initCtxCancel()

	app, err := wire.InitializeApplication(initCtx)
	if err != nil {
		log.Errorf("Failed to initialize application: %s", err)
		return
	}

	if err := app.Start(); err != nil {
		log.Errorf("Failed to start application: %s", err)
		return
	}

	stopCtx, stopCtxCancel := context.WithTimeout(context.Background(), stopTimeout*time.Second)
	defer stopCtxCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Errorf("Failed to gracefully stop application: %s", err)
		return
	}

	log.Info("Application closed")
}
