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

type TeamRepository interface {
	CreateWithMembers(ctx context.Context, teamName string, users []domain.User) (*domain.Team, error)
	GetByName(ctx context.Context, teamName string) (*domain.Team, error)
}

type UserRepository interface {
	ListByTeamID(ctx context.Context, teamID string) ([]domain.User, error)
}

type Service struct {
	lgr *slog.Logger

	teamsRepo TeamRepository
	usersRepo UserRepository
}

func New(
	lgr *slog.Logger,
	teamRepo TeamRepository,
	userRepo UserRepository,
) *Service {
	return &Service{
		lgr:       lgr,
		teamsRepo: teamRepo,
		usersRepo: userRepo,
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

	teamDB, err := s.teamsRepo.GetByName(ctx, teamName)
	if errors.Is(err, repoErr.ErrTeamNotFound) {
		lgr.DebugContext(ctx, "team not found")

		return nil, svcErr.ErrTeamNotFound
	}
	if err != nil {
		lgr.ErrorContext(ctx, "failed to get team by name", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	members, err := s.usersRepo.ListByTeamID(ctx, teamDB.ID)
	if err != nil {
		lgr.ErrorContext(ctx, "failed to list team members", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	teamDB.Members = members

	return teamDB, nil
}
