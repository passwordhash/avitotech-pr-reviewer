package user

import (
	"context"
	"fmt"

	"avitotech-pr-reviewer/internal/domain"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
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

// ListByTeamID возвращает список участников по идентификатору команды.
// Если команда не найдена или у команды нет участников, возвращается пустой слайс.
func (r *Repository) ListByTeamID(ctx context.Context, teamID string) ([]domain.Member, error) {
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

	var members []domain.Member
	for rows.Next() {
		m, err := pgPkg.RowToStructByName[model.User](rows)
		if err != nil {
			return nil, fmt.Errorf("%s: map row: %w", op, err)
		}
		members = append(members, *m.ToMemberDomain())
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return members, nil
}

// SetIsActive обновляет статус активности пользователя.
// Возвращает обновленного пользователя с именем команды.
// Если пользователь не найден, возвращается ошибка repoErr.ErrUserNotFound.
// Может вернуть repoErr.ErrTeamNotFound, если команда пользователя не найдена.
func (r *Repository) SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	const op = "repository.user.SetIsActive"

	teamName, err := r.getUsersTeamName(ctx, r.db, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: get user's team name: %w", op, err)
	}

	const updateQuery = `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING user_id, username, is_active, team_id
	`

	row := r.db.QueryRow(ctx, updateQuery, isActive, userID)
	var userDB model.User
	err = row.Scan(&userDB.UserID, &userDB.Username, &userDB.IsActive, &userDB.TeamID)
	if pgPkg.IsNoRowsError(err) {
		return nil, fmt.Errorf("%s: %w", op, repoErr.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: scan row: %w", op, err)
	}

	return userDB.ToUserDomain(teamName), nil
}

// GetByID возвращает пользователя по его идентификатору.
// Если пользователь не найден, возвращается ошибка repoErr.ErrUserNotFound.
// Может вернуть repoErr.ErrTeamNotFound, если команда пользователя не найдена.
func (r *Repository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	const op = "repository.user.GetByID"

	teamName, err := r.getUsersTeamName(ctx, r.db, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: get user's team name: %w", op, err)
	}

	const getQuery = `
		SELECT user_id, username, is_active, team_id
		FROM users
		WHERE user_id = $1
	`

	row := r.db.QueryRow(ctx, getQuery, userID)
	var userDB model.User
	err = row.Scan(&userDB.UserID, &userDB.Username, &userDB.IsActive, &userDB.TeamID)
	if pgPkg.IsNoRowsError(err) {
		return nil, fmt.Errorf("%s: %w", op, repoErr.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: scan row: %w", op, err)
	}

	return userDB.ToUserDomain(teamName), nil
}

func (r *Repository) getUsersTeamName(ctx context.Context, q pgPkg.Querier, userID string) (string, error) {
	const op = "repository.user.getUsersTeamName"

	const query = `
		SELECT t.name
		FROM users u
		JOIN teams t ON u.team_id = t.team_id
		WHERE u.user_id = $1
	`

	row := q.QueryRow(ctx, query, userID)
	var teamName string
	err := row.Scan(&teamName)
	if pgPkg.IsNoRowsError(err) {
		return "", fmt.Errorf("%s: %w", op, repoErr.ErrTeamNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("%s: scan row: %w", op, err)
	}

	return teamName, nil
}
