package migrator

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
	"log"
)

func RunMigrations(db *gorm.DB) error {
	_db, err := db.DB()
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(_db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app//migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %v", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
