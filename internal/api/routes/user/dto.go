package user

import "avitotech-pr-reviewer/internal/domain"

type User struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"isActive"`
}

func toUserFromDomain(u *domain.User) *User {
	return &User{
		ID:       u.ID,
		Username: u.Username,
		//TeamName: u.,
		IsActive: u.IsActive,
	}
}

type setIsActiveRequest struct {
	UserID   string `json:"userId" binding:"required"`
	IsActive bool   `json:"isActive" binding:"required"`
}
