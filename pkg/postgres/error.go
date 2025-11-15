package pgPkg

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const UniqueViolationCode = "23505"

func IsUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == UniqueViolationCode {
			return true
		}
	}

	return false
}

func IsNoRowsError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
