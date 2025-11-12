package main

import (
	"avitotech-pr-reviewer/internal/app"
	"avitotech-pr-reviewer/internal/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("application panic", "err", r)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	defer cancel() // TODO: есть ли смысл в этом cancel ?

	cfg := config.MustLoad()

	lgr := config.NewLogger(cfg.Env)
	slog.SetDefault(lgr)

	lgr.Info("starting pr-reviewer applicaitoin")

	application := app.New(ctx, lgr, cfg)

	go application.Srv.MustRun(ctx)

	<- ctx.Done()
}
