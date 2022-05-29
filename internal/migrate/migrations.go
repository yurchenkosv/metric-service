package migration

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(dbConnetion string) {
	m, err := migrate.New(
		"file://db/migrations",
		dbConnetion)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
