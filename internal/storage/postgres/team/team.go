package team

import (
	"context"
	"errors"
	"fmt"

	"avitotech-pr-reviewer/internal/domain"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
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

// CreateWithMembers создает команду и пользователей,
// принадлежащих к этой команде, назначает их в эту команду.
// Если переданный список пользователей содержит пользователей,
// которые уже существуют в базе, то их данные обновляются.
// Если команда с таким именем уже существует, возвращается ошибка repoErr.ErrTeamExists.
func (r *Repository) CreateWithMembers(
	ctx context.Context,
	teamName string,
	users []domain.User,
) (*domain.Team, error) {
	const op = "repository.team.CreateWithMembers"

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var team domain.Team
	const createTeamQuery = `
		INSERT INTO teams (team_name)
		VALUES ($1) RETURNING team_id, team_name
	`
	err = tx.QueryRow(ctx, createTeamQuery, teamName).Scan(&team.ID, &team.Name)
	if pgPkg.IsUniqueViolationError(err) {
		return nil, repoErr.ErrTeamExists
	}
	if err != nil {
		return nil, fmt.Errorf("%s: create team: %w", op, err)
	}

	const createTeamMemberQuery = `
		INSERT INTO users (user_id, username, is_active, team_id)
			VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id)
            DO UPDATE SET
                username = EXCLUDED.username,
                is_active = EXCLUDED.is_active,
                team_id = EXCLUDED.team_id
	`
	batch := &pgPkg.Batch{}
	for _, user := range users {
		batch.Queue(createTeamMemberQuery,
			user.ID,
			user.Username,
			user.IsActive,
			team.ID)
	}
	batchResults := tx.SendBatch(ctx, batch)
	defer func() {
		_ = batchResults.Close()
	}()

	for range users {
		_, execErr := batchResults.Exec()
		if execErr != nil {
			e := batchResults.Close()
			err = errors.Join(e, fmt.Errorf("%s: %w", op, execErr))

			return nil, err
		}
	}

	err = batchResults.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	team.Members = users

	return &team, nil
}
