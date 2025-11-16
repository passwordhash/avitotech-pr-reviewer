package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

func (s PRStatus) IsValid() bool {
	switch s {
	case PRStatusOpen, PRStatusMerged:
		return true
	default:
		return false
	}
}

type PullRequest struct {
	ID                  string
	Name                string
	AuthorID            string
	InNeedMoreReviewers bool
	Status              PRStatus
	Reviewers           []string
	CreatedAt           time.Time
	MergedAt            *time.Time
}
