package model

import "avitotech-pr-reviewer/internal/domain"

type User struct {
	UserID   string `db:"user_id"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active"`
	TeamID   string `db:"team_id"`
}

func (u User) ToUserDomain() *domain.User {
	return &domain.User{
		ID:       u.UserID,
		Username: u.Username,
		IsActive: u.IsActive,
		TeamID:   u.TeamID,
	}
}

func (u User) ToMemberDomain() *domain.Member {
	return &domain.Member{
		ID:       u.UserID,
		Username: u.Username,
		IsActive: u.IsActive,
	}
}
