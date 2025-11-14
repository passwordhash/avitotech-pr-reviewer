// Package app содержит инициализацию всех составляющих приложения,
// например, HTTP сервера, баз данных и т.д.
package app

import (
	"context"
	"log/slog"

	httpapp "avitotech-pr-reviewer/internal/app/http"
	"avitotech-pr-reviewer/internal/config"
	teamService "avitotech-pr-reviewer/internal/service/team"
)

type App struct {
	Srv *httpapp.App
}

func New(
	_ context.Context,
	lgr *slog.Logger,
	cfg *config.Config,
) *App {
	// pkgPool ...

	// repos ...
	teamSvc := teamService.New(lgr.WithGroup("team_service"))

	srv := httpapp.New(
		lgr.WithGroup("httpapp"),
		teamSvc,
		httpapp.WithPort(cfg.HTTP.Port),
		httpapp.WithReadTimeout(cfg.HTTP.ReadTimeout),
		httpapp.WithWriteTimeout(cfg.HTTP.ReadTimeout),
		httpapp.WithRequestTimeout(cfg.HTTP.ReadTimeout),
	)

	return &App{
		Srv: srv,
	}
}
