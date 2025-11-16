package model

import (
	"database/sql"
	"time"

	"avitotech-pr-reviewer/internal/domain"
)

type PullRequest struct {
	ID                  string       `db:"pull_request_id"`
	Name                string       `db:"pull_request_name"`
	AuthorID            string       `db:"author_id"`
	StatusID            string       `db:"status_id"`
	InNeedMoreReviewers bool         `db:"is_need_more_reviewers"`
	CreatedAt           time.Time    `db:"created_at"`
	MergedAt            sql.NullTime `db:"merged_at"`
}

func (pr *PullRequest) ToDomain(status domain.PRStatus) *domain.PullRequest {
	var mergedAt *time.Time
	if pr.MergedAt.Valid {
		mergedAt = &pr.MergedAt.Time
	}

	return &domain.PullRequest{
		ID:                  pr.ID,
		Name:                pr.Name,
		AuthorID:            pr.AuthorID,
		Status:              status,
		InNeedMoreReviewers: pr.InNeedMoreReviewers,
		CreatedAt:           pr.CreatedAt,
		MergedAt:            mergedAt,
	}
}
