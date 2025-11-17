package pullrequest

import (
	"time"

	"avitotech-pr-reviewer/internal/domain"
)

type PullRequest struct {
	ID                  string     `json:"pull_request_id"`
	Name                string     `json:"pull_request_name"`
	AuthorID            string     `json:"author_id"`
	Status              string     `json:"status"`
	AssignedReviewerIDs []string   `json:"assigned_reviewers"`
	MergedAt            *time.Time `json:"merged_at,omitempty"`
}

func FromDomainPR(pr *domain.PullRequest) *PullRequest {
	return &PullRequest{
		ID:                  pr.ID,
		Name:                pr.Name,
		AuthorID:            pr.AuthorID,
		Status:              string(pr.Status),
		AssignedReviewerIDs: pr.Reviewers,
		MergedAt:            pr.MergedAt,
	}
}

type CreatePullRequestRequest struct {
	ID       string `json:"pull_request_id" binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}
