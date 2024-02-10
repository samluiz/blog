package pkg

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/crypto/bcrypt"
)

func connect() (*sqlx.DB, error) {
	url := os.Getenv("DATABASE_URL")

	db, err := sqlx.Connect("libsql", url)

	if err != nil {
		return nil, err
	}

	err = initTables()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initTables() error {
	db, err := connect()
	if err != nil {
		return err
	}

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return err
	}

	username := os.Getenv("ADMIN_USERNAME")
	hashed, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), 8)
	if err != nil {
		return err
	}

	password := string(hashed)

	_, err = db.Exec("INSERT INTO users (username, password, is_admin) VALUES ($1, $2, $3)", username, password, 1)

	return nil
}