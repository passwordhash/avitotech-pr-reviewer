package team

import (
	"log/slog"

	"avitotech-pr-reviewer/internal/domain"
)

type Service struct {
	lgr *slog.Logger
}

func New(lgr *slog.Logger) *Service {
	return &Service{lgr: lgr}
}

func (s *Service) Add(teamName string, members []domain.User) error {
	panic("implement me")
}
