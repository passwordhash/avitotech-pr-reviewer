package user

import (
	"context"
	"fmt"

	"avitotech-pr-reviewer/internal/domain"
	"avitotech-pr-reviewer/internal/storage/postgres/user/model"
	pgPkg "avitotech-pr-reviewer/pkg/postgres"
)

type Repository struct {
	db pgPkg.DB
}

func New(db pgPkg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// ListByTeamID возвращает список пользователей по идентификатору команды.
// Если команда не найдена или у команды нет участников, возвращается пустой слайс.
func (r *Repository) ListByTeamID(
	ctx context.Context,
	teamID string,
) ([]domain.User, error) {
	const op = "repository.team.ListByTeamID"

	const listQuery = `
		SELECT user_id, username, is_active, team_id
		FROM users
		WHERE team_id = $1
	`

	rows, err := r.db.Query(ctx, listQuery, teamID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		usr, err := pgPkg.RowToStructByName[model.User](rows)
		if err != nil {
			return nil, fmt.Errorf("%s: map row: %w", op, err)
		}
		users = append(users, *usr.ToDomain())
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}
