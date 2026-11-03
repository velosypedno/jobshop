package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/velosypedno/jobshop/internal/backend"
	"go.uber.org/zap"
)

func main() {
	l := backend.NewLogger()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	cfg, err := backend.Load()
	if err != nil {
		l.Panic("Failed to load config", zap.Error(err))
	}
	l.Info("Creating app...")
	app := backend.New(cfg, l)
	l.Info("App created")

	l.Info("Running app...")
	if err := app.Run(ctx); err != nil {
		l.Panic("Failed to run app", zap.Error(err))
	}

	l.Info("App stopped successfully")
}
