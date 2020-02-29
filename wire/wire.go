//+build wireinject

package wire

import (
	"context"
	"microservice/internal/app/domain"
	"microservice/internal/pkg/mongodb"

	"microservice/internal/app"
	"microservice/internal/app/drivers/rest"
	"microservice/internal/pkg/viper"

	"github.com/google/wire"
)

func InitializeApplication(ctx context.Context) (app.App, error) {
	wire.Build(
		viper.NewConfiguration,
		wire.Bind(new(rest.Configuration), new(*viper.Service)),
		wire.Bind(new(mongodb.Configuration), new(*viper.Service)),

		domain.NewDomain,
		wire.Bind(new(rest.DomainSvc), new(*domain.Domain)),

		mongodb.NewClient,
		wire.Bind(new(domain.DocumentDB), new(*mongodb.MongoDB)),

		rest.NewServer,
		wire.Bind(new(app.RestServer), new(*rest.Server)),

		app.NewApp,
	)
	return app.App{}, nil
}
