package model

import "avitotech-pr-reviewer/internal/domain"

type User struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
	TeamID   string `db:"team_id"`
}

func (u User) ToDomain() *domain.User {
	return &domain.User{
		Username: u.Username,
		IsActive: u.IsActive,
	}
}
