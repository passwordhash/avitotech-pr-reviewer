package repoErr

import "errors"

var (
	ErrTeamExists = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
)
