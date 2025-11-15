// Package app содержит инициализацию всех составляющих приложения,
// например, HTTP сервера, баз данных и т.д.
package app

import (
	"context"
	"log/slog"

	httpapp "avitotech-pr-reviewer/internal/app/http"
	"avitotech-pr-reviewer/internal/config"
	teamService "avitotech-pr-reviewer/internal/service/team"
	userService "avitotech-pr-reviewer/internal/service/user"
	teamRepository "avitotech-pr-reviewer/internal/storage/postgres/team"
	userRepository "avitotech-pr-reviewer/internal/storage/postgres/user"
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
	userRepo := userRepository.New(pgPool)

	teamSvc := teamService.New(lgr.WithGroup("service.team"), teamRepo, userRepo)
	userSvc := userService.New(lgr.WithGroup("service.user"), userRepo, teamRepo)

	srv := httpapp.New(
		lgr,
		teamSvc,
		userSvc,
		httpapp.WithPort(cfg.HTTP.Port),
		httpapp.WithReadTimeout(cfg.HTTP.ReadTimeout),
		httpapp.WithWriteTimeout(cfg.HTTP.ReadTimeout),
		httpapp.WithRequestTimeout(cfg.HTTP.ReadTimeout),
	)

	return &App{
		Srv: srv,
	}
}
