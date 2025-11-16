// Package app содержит инициализацию всех составляющих приложения,
// например, HTTP сервера, баз данных и т.д.
package app

import (
	"context"
	"log/slog"

	httpapp "avitotech-pr-reviewer/internal/app/http"
	"avitotech-pr-reviewer/internal/config"
	prService "avitotech-pr-reviewer/internal/service/pullrequest"
	teamService "avitotech-pr-reviewer/internal/service/team"
	userService "avitotech-pr-reviewer/internal/service/user"
	prRepository "avitotech-pr-reviewer/internal/storage/postgres/pullrequest"
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
	prRepo := prRepository.New(pgPool)

	teamSvc := teamService.New(lgr.WithGroup("service.team"), teamRepo, userRepo)
	userSvc := userService.New(lgr.WithGroup("service.user"), userRepo, teamRepo, cfg.App.AdminToken)
	prSvc := prService.New(lgr.WithGroup("service.pullrequest"),
		prRepo, userRepo, teamRepo, cfg.App.MaxReviewersPerPR)

	srv := httpapp.New(
		lgr,
		teamSvc,
		userSvc,
		prSvc,
		httpapp.WithPort(cfg.HTTP.Port),
		httpapp.WithReadTimeout(cfg.HTTP.ReadTimeout),
		httpapp.WithWriteTimeout(cfg.HTTP.WriteTimeout),
		httpapp.WithRequestTimeout(cfg.HTTP.GatewayTimeout),
	)

	return &App{
		Srv: srv,
	}
}
