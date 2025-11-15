// Package app содержит инициализацию всех составляющих приложения,
// например, HTTP сервера, баз данных и т.д.
package app

import (
	"context"
	"log/slog"

	httpapp "avitotech-pr-reviewer/internal/app/http"
	"avitotech-pr-reviewer/internal/config"
	teamService "avitotech-pr-reviewer/internal/service/team"
	teamRepository "avitotech-pr-reviewer/internal/storage/postgres/team"
	pgPkg "avitotech-pr-reviewer/pkg/postgres"
)

type App struct {
	Srv *httpapp.App
}

func New(
	ctx context.Context,
	lgr *slog.Logger,
	cfg *config.Config,
) *App {
	pgPool, err := pgPkg.NewPool(ctx, cfg.PG.DSN(), pgPkg.WithMaxConns(cfg.PG.MaxConns))
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	teamRepo := teamRepository.New(pgPool)

	teamSvc := teamService.New(lgr.WithGroup("team_service"), teamRepo)

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
