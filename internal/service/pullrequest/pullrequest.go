package pullrequest

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"

	"avitotech-pr-reviewer/internal/domain"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
)

type prRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
}

type userRepository interface {
	GetByID(ctx context.Context, userID string) (*domain.User, error)
}

type teamRepository interface {
	GetActiveMembersByTeamID(ctx context.Context, teamID string) ([]domain.Member, error)
}

type Service struct {
	lgr *slog.Logger

	prRepo   prRepository
	userRepo userRepository
	teamRepo teamRepository

	maxReviewers int // максимальное количество ревьюверов на PR
}

func New(
	lgr *slog.Logger,
	prRepo prRepository,
	userRepo userRepository,
	teamRepo teamRepository,
	maxReviewers int,
) *Service {
	return &Service{
		lgr:          lgr,
		prRepo:       prRepo,
		userRepo:     userRepo,
		teamRepo:     teamRepo,
		maxReviewers: maxReviewers,
	}
}

// CreatePullRequest создаёт новый Pull Request с указанным ID, именем и автором.
// Выбирает ревьюверов из активных участников команды автора.
// Если Pull Request с таким ID уже существует, возвращается ошибка svcErr.ErrPRExists.
// Если автор не найден, возвращается ошибка svcErr.ErrUserNotFound.
func (s *Service) CreatePullRequest(ctx context.Context, id, name, authorID string) (*domain.PullRequest, error) {
	const op = "pullrequest.CreatePullRequest"

	lgr := s.lgr.With(
		slog.String("op", op),
		slog.String("pull_request_id", id),
		slog.String("name", name),
		slog.String("author_id", authorID),
	)

	author, err := s.userRepo.GetByID(ctx, authorID)
	if errors.Is(err, repoErr.ErrUserNotFound) {
		lgr.DebugContext(ctx, "author not found", slog.String("error", err.Error()))

		return nil, svcErr.ErrUserNotFound
	}
	if err != nil {
		lgr.ErrorContext(ctx, "failed to get author by ID", slog.String("error", err.Error()))

		return nil, err
	}

	teamMembers, err := s.teamRepo.GetActiveMembersByTeamID(ctx, author.TeamID)
	if err != nil {
		lgr.ErrorContext(ctx, "failed to get team members by team ID", slog.String("error", err.Error()))

		return nil, err
	}

	reviewers := s.selectReviewers(teamMembers, authorID, s.maxReviewers)

	pr := &domain.PullRequest{
		ID:                  id,
		Name:                name,
		AuthorID:            authorID,
		InNeedMoreReviewers: len(reviewers) < 1,
		Status:              domain.PRStatusOpen,
		Reviewers:           reviewers,
	}

	pr, err = s.prRepo.Create(ctx, pr)
	if errors.Is(err, repoErr.ErrPRExists) {
		lgr.DebugContext(ctx, "pull request already exists", slog.String("error", err.Error()))

		return nil, svcErr.ErrPRExists
	}
	if err != nil {
		lgr.ErrorContext(ctx, "failed to create pull request", slog.String("error", err.Error()))

		return nil, err
	}

	return pr, nil
}

func (s *Service) selectReviewers(teamMembers []domain.Member, authorID string, maxCount int) []string {
	candidates := make([]string, 0, len(teamMembers))
	for _, member := range teamMembers {
		if member.ID != authorID {
			candidates = append(candidates, member.ID)
		}
	}

	if len(candidates) <= maxCount {
		return candidates
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	return candidates[:maxCount]
}
