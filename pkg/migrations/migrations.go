package migrations

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(dsn string, migrationsFolderPath string) error {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	driver, err := pgx.WithInstance(sqlDB, &pgx.Config{})
	if err != nil {
		return err
	}

	dbName, err := dbNameByDSN(dsn)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsFolderPath, dbName, driver)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func dbNameByDSN(dsn string) (string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	if parsedURL.Scheme != "postgres" {
		return "", ErrIncorrectDatabaseSchema
	}

	dbName := strings.TrimPrefix(parsedURL.Path, "/")
	if dbName == "" {
		return "", ErrNoSpecifiedDatabaseName
	}

	return dbName, nil
}
