package team

import (
	"context"
	"errors"
	"fmt"

	"avitotech-pr-reviewer/internal/domain"
	repoErr "avitotech-pr-reviewer/internal/storage/errors"
	"avitotech-pr-reviewer/internal/storage/postgres/team/model"
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
	members []domain.Member,
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
		INSERT INTO user (user_id, username, is_active, team_id)
			VALUES ($1, $2, $3, $4)
            ON CONFLICT (user_id)
            DO UPDATE SET
                username = EXCLUDED.username,
                is_active = EXCLUDED.is_active,
                team_id = EXCLUDED.team_id
	`
	batch := &pgPkg.Batch{}
	for _, member := range members {
		batch.Queue(createTeamMemberQuery,
			member.ID,
			member.Username,
			member.IsActive,
			team.ID)
	}
	batchResults := tx.SendBatch(ctx, batch)
	defer func() {
		_ = batchResults.Close()
	}()

	for range members {
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

	team.Members = members

	return &team, nil
}

// GetByName возвращает команду по ее имени.
// Если команда с таким именем не найдена, возвращается ошибка repoErr.ErrTeamNotFound.
func (r *Repository) GetByName(ctx context.Context, teamName string) (*domain.Team, error) {
	const op = "repository.team.GetByName"

	const getQuery = `
		SELECT * FROM teams
		WHERE team_name = $1
	`
	rows, err := r.db.Query(ctx, getQuery, teamName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	teamDB, err := pgPkg.CollectExactlyOneRow(rows, pgPkg.RowToStructByName[model.Team])
	if pgPkg.IsNoRowsError(err) {
		return nil, repoErr.ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return teamDB.ToDomain(), nil
}

// GetByID возвращает команду по ее идентификатору.
// Если команда с таким идентификатором не найдена, возвращается ошибка repoErr.ErrTeamNotFound.
func (r *Repository) GetByID(ctx context.Context, teamID string) (*domain.Team, error) {
	const op = "repository.team.GetByID"

	const getQuery = `
		SELECT * FROM teams
		WHERE team_id = $1
	`
	rows, err := r.db.Query(ctx, getQuery, teamID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	teamDB, err := pgPkg.CollectExactlyOneRow(rows, pgPkg.RowToStructByName[model.Team])
	if pgPkg.IsNoRowsError(err) {
		return nil, repoErr.ErrTeamNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return teamDB.ToDomain(), nil
}
