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
	log.Printf("Applied %d migrations!\n", n)
	return nil
}

func IsMigrated(m *migrate.FileMigrationSource, db *sql.DB) (bool, error) {
	pending, err := m.FindMigrations()
	if err != nil {
		return false, err
	}
	applied, err := migrate.GetMigrationRecords(db, "postgres")
	if err != nil {
		return false, err
	}
	if len(pending) == len(applied) {
		return true, nil
	}
	return false, nil
}


