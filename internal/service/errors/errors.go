package svcErr

import "errors"

var (
	ErrTeamExists = errors.New("team already exists")
)
