package pullrequest

import "avitotech-pr-reviewer/internal/domain"

type PullRequest struct {
	ID                 string
	Name               string
	AuthorID           string
	Status             string
	AssignedReviewerID []string
}

func FromDomainPR(pr *domain.PullRequest) *PullRequest {
	return &PullRequest{
		ID:       pr.ID,
		Name:     pr.Name,
		AuthorID: pr.AuthorID,
		Status: string(pr.Status),
		AssignedReviewerID: pr.Reviewers,
	}
}

type CreatePullRequestRequest struct {
	ID       string `json:"pull_request_id" binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}
