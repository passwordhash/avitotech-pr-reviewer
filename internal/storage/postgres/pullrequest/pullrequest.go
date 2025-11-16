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

func (r *Repository) addReviewers(ctx context.Context, q pgPkg.Querier, prID string, reviewerIDs []string) error {
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

func (r *Repository) getStatusByID(
	ctx context.Context,
	q pgPkg.Querier,
	statusID string,
) (domain.PRStatus, error) {
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
