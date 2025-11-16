package pullrequest

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"avitotech-pr-reviewer/internal/domain"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
	"avitotech-pr-reviewer/internal/storage/postgres/pullrequest/model"
	pgPkg "avitotech-pr-reviewer/pkg/postgres"
)

type Repository struct {
	db pgPkg.DB
}

func New(db pgPkg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create создаёт новый Pull Request.
// Если Pull Request с таким ID уже существует, возвращается ошибка repoErr.ErrPRExists.
func (r *Repository) Create(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	const op = "pullrequest.Repository.Create"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const query = `
        INSERT INTO pull_requests (
            pull_request_id,
            pull_request_name,
            author_id
        )
        VALUES (@id, @name, @author_id)
        RETURNING pull_request_id, pull_request_name, author_id,
				  created_at, status_id, merged_at, is_need_more_reviewers
    `

	rows, err := tx.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"id":        pr.ID,
			"name":      pr.Name,
			"author_id": pr.AuthorID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.PullRequest])
	if pgPkg.IsUniqueViolationError(err) {
		return nil, repoErr.ErrPRExists
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	status, err := r.getStatusByID(ctx, tx, created.StatusID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pullRequest := created.ToDomain(status)

	if len(pr.Reviewers) > 0 {
		err = r.addReviewers(ctx, tx, pr.ID, pr.Reviewers)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	pullRequest.Reviewers = pr.Reviewers

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pullRequest, nil
}

// GetByID возвращает Pull Request по его ID.
// Если Pull Request не найден, возвращается ошибка repoErr.ErrPRNotFound.
func (r *Repository) GetByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	const op = "pullrequest.Repository.GetByID"

	const query = `
		SELECT pull_request_id, pull_request_name, author_id,
			   created_at, status_id, merged_at, is_need_more_reviewers
		FROM pull_requests
		WHERE pull_request_id = $1
	`
	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	found, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.PullRequest])
	if pgPkg.IsNoRowsError(err) {
		return nil, repoErr.ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	status, err := r.getStatusByID(ctx, r.db, found.StatusID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return found.ToDomain(status), nil
}

// GetReviewerIDs возвращает список ID ревьюеров, назначенных на указанный Pull Request.
// Метод не возвращает ошибку, если Pull Request не найден или у него нет назначенных ревьюеров.
func (r *Repository) GetReviewerIDs(ctx context.Context, prID string) ([]string, error) {
	const op = "pullrequest.Repository.GetReviewerIDs"

	const query = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
	`
	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var reviewerIDs []string
	for rows.Next() {
		var reviewerID string
		err := rows.Scan(&reviewerID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviewerIDs = append(reviewerIDs, reviewerID)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return reviewerIDs, nil
}

// SetMerged помечает указанный Pull Request как merged.
// Если Pull Request не найден, возвращается ошибка repoErr.ErrPRNotFound.
// Операция не является идемпотентной. Нужно вызывать только если Pull Request ещё не был помечен как merged.
func (r *Repository) SetMerged(ctx context.Context, prID string) (*domain.PullRequest, error) {
	const op = "pullrequest.Repository.SetMerged"

	const query = `
		UPDATE pull_requests
		SET status_id = (SELECT id FROM pull_request_statuses WHERE UPPER(status) = 'MERGED'),
			merged_at = NOW()
		WHERE pull_request_id = $1
		RETURNING *
	`
	rows, err := r.db.Query(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.PullRequest])
	if pgPkg.IsNoRowsError(err) {
		return nil, repoErr.ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	status, err := r.getStatusByID(ctx, r.db, updated.StatusID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pullRequest := updated.ToDomain(status)

	return pullRequest, nil
}

func (r *Repository) addReviewers(ctx context.Context, q pgPkg.Tx, prID string, reviewerIDs []string) error {
	const op = "pullrequest.Repository.addReviewers"

	const query = `
		INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`

	batch := &pgx.Batch{}
	for _, reviewerID := range reviewerIDs {
		batch.Queue(query, prID, reviewerID)
	}
	batchResults := q.SendBatch(ctx, batch)
	defer func() {
		_ = batchResults.Close()
	}()

	for range reviewerIDs {
		_, execErr := batchResults.Exec()
		if execErr != nil {
			_ = batchResults.Close()

			if pgPkg.IsForeignKeyErr(execErr) {
				return fmt.Errorf("%s: %w", op, repoErr.ErrUserNotFound)
			}

			return fmt.Errorf("%s: %w", op, execErr)
		}
	}

	return nil
}

func (r *Repository) getStatusByID(ctx context.Context, q pgPkg.Querier, statusID string) (domain.PRStatus, error) {
	const op = "pullrequest.Repository.getStatusByID"
	const query = `
        SELECT status
        FROM pull_request_statuses
        WHERE id = $1
    `
	var status string
	err := q.QueryRow(ctx, query, statusID).Scan(&status)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	domainStatus := domain.PRStatus(status)
	if !domainStatus.IsValid() {
		return "", repoErr.ErrInvalidStatus
	}

	return domainStatus, nil
}
