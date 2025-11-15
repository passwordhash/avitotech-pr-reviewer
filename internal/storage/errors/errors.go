package repoErr

import "errors"

var (
	ErrTeamExists = errors.New("team already exists")
)
