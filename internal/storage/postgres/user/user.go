package user

import pgPkg "avitotech-pr-reviewer/pkg/postgres"

type Repository struct {
	db pgPkg.DB
}

func New(db pgPkg.DB) *Repository {
	return &Repository{
		db: db,
	}
}
