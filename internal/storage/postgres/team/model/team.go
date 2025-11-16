package model

import "avitotech-pr-reviewer/internal/domain"

type Team struct {
	ID   string `db:"team_id"`
	Name string `db:"team_name"`
}

func (t Team) ToDomain() *domain.Team {
	return &domain.Team{
		ID:   t.ID,
		Name: t.Name,
	}
}
