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
	TeamWithMembers(ctx context.Context, teamName string) (*domain.TeamNormalized, error)
}

type Service struct {
	lgr *slog.Logger

	teamsRepo TeamsRepository
}

func New(
	lgr *slog.Logger,
	reamsRepo TeamsRepository,
) *Service {
	return &Service{
		teamsRepo: reamsRepo,
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

	lgr := s.lgr.With(
		slog.String("op", op),
		slog.String("teamName", teamName),
	)

	createdTeam, err := s.teamsRepo.CreateWithMembers(ctx, teamName, members)
	if errors.Is(err, repoErr.ErrTeamExists) {
		lgr.DebugContext(ctx, "team already exists")

		return nil, svcErr.ErrTeamExists
	}
	if err != nil {
		lgr.ErrorContext(ctx, "failed to create team with members", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lgr.InfoContext(ctx, "team with members created successfully",
		slog.String("teamID", createdTeam.ID),
		slog.Int("membersCount", len(createdTeam.Members)),
	)

	return createdTeam, nil
}

// TeamWithMembers возвращает команду с указанным именем вместе с ее участниками.
func (s *Service) TeamWithMembers(
	ctx context.Context,
	teamName string,
) (*domain.Team, error) {
	const op = "team.TeamWithMembers"

	lgr := s.lgr.With(
		slog.String("op", op),
		slog.String("teamName", teamName),
	)
	_ = lgr

	return nil, nil //nolint:nilnil
}
