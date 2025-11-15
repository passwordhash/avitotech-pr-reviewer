package team

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"avitotech-pr-reviewer/internal/domain"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
)

type TeamsRepository interface {
	CreateWithMembers(ctx context.Context, teamName string, users []domain.User) (*domain.Team, error)
}

type Service struct {
	lgr *slog.Logger

	reamsRepo TeamsRepository
}

func New(
	lgr *slog.Logger,
	reamsRepo TeamsRepository,
) *Service {
	return &Service{
		reamsRepo: reamsRepo,
		lgr:       lgr,
	}
}

// CreateTeam создает команду с указанным именем и участниками.
// Если команда с таким именем уже существует, возвращается ошибка svcErr.ErrTeamExists.
func (s *Service) CreateTeam(
	ctx context.Context,
	teamName string,
	members []domain.User,
) (*domain.Team, error) {
	const op = "team.CreateTeam"

	lgr := s.lgr.With("op", op, "teamName", teamName)

	createdTeam, err := s.reamsRepo.CreateWithMembers(ctx, teamName, members)
	if errors.Is(err, repoErr.ErrTeamExists) {
		lgr.DebugContext(ctx, "team already exists")

		return nil, svcErr.ErrTeamExists
	}
	if err != nil {
		lgr.ErrorContext(ctx, "failed to create team with members", "err", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lgr.InfoContext(ctx, "team with members created successfully",
		slog.String("teamID", createdTeam.ID),
		slog.Int("membersCount", len(createdTeam.Members)))

	return createdTeam, nil
}
