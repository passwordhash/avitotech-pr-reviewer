package model

type Team struct {
	ID   string `db:"team_id"`
	Name string `db:"team_name"`
}
