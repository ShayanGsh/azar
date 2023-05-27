package db

import (
	"database/sql"
	"log"

	"github.com/rubenv/sql-migrate"
)

func Migration(path string) *migrate.FileMigrationSource {
	return &migrate.FileMigrationSource{
		Dir: path,
	}
}

func RunMigration(m *migrate.FileMigrationSource, db *sql.DB, d migrate.MigrationDirection) error {
	n, err := migrate.Exec(db, "postgres", m, d)
	if err != nil {
		return err
	}
	log.Printf("Applied %d migrations", n)
	return nil
}

func IsMigrated(m *migrate.FileMigrationSource, db *sql.DB) (bool, error) {
	records, err := m.FindMigrations()

	if err != nil {
		return false, err
	}
	if len(records) == 0 {
		return false, nil
	}
	n, err := migrate.ExecMax(db, "postgres", m, migrate.Up, len(records))
	if err != nil {
		return false, err
	}
	return n == len(records), nil
}


