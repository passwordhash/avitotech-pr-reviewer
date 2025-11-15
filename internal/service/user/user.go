package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"avitotech-pr-reviewer/internal/domain"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
)

type UserRepository interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
}

type TeamRepository interface {
	GetByID(ctx context.Context, teamID string) (*domain.Team, error)
}

type Service struct {
	lgr *slog.Logger

	userRepo UserRepository
	teamRepo TeamRepository
}

func New(
	lgr *slog.Logger,
	userRepo UserRepository,
	teamRepo TeamRepository,
) *Service {
	return &Service{
		lgr:      lgr,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

// SetIsActive ...
func (s Service) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	const op = "user.SetIsActive"

	lgr := s.lgr.With(
		slog.String("op", op),
		slog.String("userID", userID),
		slog.Bool("isActive", isActive),
	)

	userItem, err := s.userRepo.SetIsActive(ctx, userID, isActive)
	if errors.Is(err, repoErr.ErrUserNotFound) {
		lgr.DebugContext(ctx, "user not found", slog.Any("error", err))

		return nil, svcErr.ErrUserNotFound
	}
	if err != nil {
		lgr.ErrorContext(ctx, "failed to set is active", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//nolint:nolintlint,godox    // TODO: операция получается не атомарной, нужно подумать как это исправить
	lgr.Info("user active status updated successfully")

	teamItem, err := s.teamRepo.GetByID(ctx, userItem.TeamID)
	if errors.Is(err, repoErr.ErrTeamNotFound) {
		lgr.DebugContext(ctx, "team not found", slog.Any("error", err))

		return nil, svcErr.ErrTeamNotFound
	}
	if err != nil {
		lgr.ErrorContext(ctx, "failed to set is active", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userItem.TeamName = teamItem.Name

	return userItem, nil
}
