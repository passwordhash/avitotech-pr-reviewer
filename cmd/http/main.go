package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avitotech-pr-reviewer/internal/app"
	"avitotech-pr-reviewer/internal/config"
)

const (
	shutdownTimeout = 10 * time.Second
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("application panic", "err", r)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := config.MustLoad()

	lgr := config.NewLogger(cfg.App.Env)
	slog.SetDefault(lgr)

	lgr.Info("starting pr-reviewer application")

	application := app.New(ctx, lgr, cfg)

	go application.Srv.MustRun(ctx)

	<-ctx.Done()

	lgr.Info("received stop signal")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := application.Srv.Stop(shutdownCtx)
	if err != nil {
		lgr.Error("failed to stop http_server gracefully", "err", err)
	} else {
		lgr.Info("PR Reviewer application http_server stopped gracefully")
	}
}
