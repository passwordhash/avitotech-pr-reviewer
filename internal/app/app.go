// Package app содержит инициализацию всех составляющих приложения,
// например, HTTP сервера, баз данных и т.д.
package app

import (
	"context"
	"log/slog"

	httpapp "avitotech-pr-reviewer/internal/app/http"
	"avitotech-pr-reviewer/internal/config"
)

type App struct {
	Srv *httpapp.App
}

func New(
	_ context.Context,
	lgr *slog.Logger,
	cfg *config.Config,
) *App {
	srv := httpapp.New(
		lgr.WithGroup("http"),
		httpapp.WithPort(cfg.HTTP.Port),
		httpapp.WithReadTimeout(cfg.HTTP.ReadTimeout),
		httpapp.WithWriteTimeout(cfg.HTTP.ReadTimeout),
		httpapp.WithRequestTimeout(cfg.HTTP.ReadTimeout),
	)

	return &App{
		Srv: srv,
	}
}
