package migrations

import "errors"

var (
	ErrIncorrectDatabaseSchema = errors.New("incorrect database schema")
	ErrNoSpecifiedDatabaseName = errors.New("no database name provided")
)
