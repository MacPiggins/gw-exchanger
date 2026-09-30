package main

import (
	"log/slog"
	"os"

	"github.com/MacPiggins/gw-exchanger/internal/app"
	"github.com/MacPiggins/gw-exchanger/internal/config"
	"github.com/MacPiggins/gw-exchanger/internal/logging"
)

func main() {
	h := &logging.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
	slog.SetDefault(slog.New(h))

	conf, err := config.Load("config.env")
	if err != nil {
		slog.Error("error while loading config", slog.Any("error", err))
		os.Exit(1)
	}
	if err := app.Run(conf); err != nil {
		slog.Error("application stopped", slog.Any("error", err))
		os.Exit(1)
	}
}
