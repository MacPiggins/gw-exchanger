package app

import (
	"context"
	"log/slog"
	"net"

	"github.com/MacPiggins/gw-exchanger/internal/config"
	"github.com/MacPiggins/gw-exchanger/internal/server"
	"github.com/MacPiggins/gw-exchanger/internal/service"
	"github.com/MacPiggins/gw-exchanger/internal/storage/postgres"
)

func Run(conf *config.Config) error {
	ctx := context.Background()
	storage, err := postgres.New(ctx, conf.Connstr)
	if err != nil {
		slog.Error("error while creating storage", slog.Any("error", err))
		return err
	}
	defer storage.Close()

	err = storage.AutoMigrate(ctx)
	if err != nil {
		slog.Error("error while migrating storage", slog.Any("error", err))
		return err
	}

	exchangeService := service.NewExchangeService(storage)
	server := server.New(exchangeService)

	slog.Info("starting exchanger server", slog.String("address", conf.Address))
	lis, err := net.Listen("tcp4", conf.Address)
	if err != nil {
		slog.Error("error while listening to address", slog.Any("error", err))
		return err
	}
	return server.Run(lis)
}
