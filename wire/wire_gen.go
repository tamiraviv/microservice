// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"context"
	"microservice/internal/app"
	"microservice/internal/app/domain"
	"microservice/internal/app/drivers/rest"
	"microservice/internal/pkg/mongodb"
	"microservice/internal/pkg/viper"
)

// Injectors from wire.go:

func InitializeApplication(ctx context.Context) (app.App, error) {
	service, err := viper.NewConfiguration()
	if err != nil {
		return app.App{}, err
	}
	mongoDB, err := mongodb.NewClient(ctx, service)
	if err != nil {
		return app.App{}, err
	}
	domainDomain, err := domain.NewDomain(mongoDB)
	if err != nil {
		return app.App{}, err
	}
	server, err := rest.NewServer(service, domainDomain)
	if err != nil {
		return app.App{}, err
	}
	appApp := app.NewApp(server)
	return appApp, nil
}