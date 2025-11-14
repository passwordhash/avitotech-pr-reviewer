package team

import "avitotech-pr-reviewer/internal/domain"

type userReq struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (u *userReq) ToDomain() domain.User {
	return domain.User{
		ID:       u.UserID,
		Username: u.Username,
		IsActive: u.IsActive,
	}
}

type addReq struct {
	TeamName string    `json:"team_name" binding:"required"`
	Members  []userReq `json:"members"`
}

func (a *addReq) ToDomainMembers() []domain.User {
	domainMembers := make([]domain.User, len(a.Members))
	for i, member := range a.Members {
		domainMembers[i] = member.ToDomain()
	}

	return domainMembers
}
