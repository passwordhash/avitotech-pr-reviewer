package team

import (
	"context"
	"log/slog"

	"avitotech-pr-reviewer/internal/domain"
)

type Service struct {
	lgr *slog.Logger
}

func New(lgr *slog.Logger) *Service {
	return &Service{lgr: lgr}
}

// CreateTeam создает команду с указанным именем и участниками.
// ...
func (s Service) CreateTeam(ctx context.Context, teamName string, members []domain.User) error {
	const op = "team.CreateTeam"

	lgr := s.lgr.With("op", op, "teamName", teamName)
	_ = lgr

	panic("implement me")
}
