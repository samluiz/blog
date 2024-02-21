package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/crypto/bcrypt"
)

func NewConnection() (*sqlx.DB, error) {
	url := os.Getenv("DATABASE_URL")

	db, err := sqlx.Connect("libsql", url)

	if err != nil {
		return nil, err
	}

	err = initTables(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadSchema() ([]byte, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(dir, "schema.sql")

	schema, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func initUser(db *sqlx.DB) error {
	log.Default().Println("Initializing admin user...")
	username := os.Getenv("ADMIN_USERNAME")

	var userExists bool

	err := db.Get(&userExists, "SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)", username)

	if err != nil {
		log.Default().Printf("Error checking if user exists: %v", err)
		return err
	}

	if userExists {
		log.Default().Println("Admin user already exists")
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), 8)
	if err != nil {
		log.Default().Printf("Error hashing password: %v", err)
		return err
	}

	password := string(hashed)

	_, err = db.Exec("INSERT INTO users (username, password, is_admin) VALUES (?, ?, ?)", username, password, 1)

	if err != nil {
		log.Default().Printf("Error creating admin user: %v", err)
		return err
	}
	return nil
}

func initTables(db *sqlx.DB) error {
	log.Default().Println("Initializing tables...")
	schema, err := loadSchema()
	if err != nil {
		log.Default().Printf("Error loading schema: %v", err)
		return err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Default().Printf("Error creating tables: %v", err)
		return err
	}

	err = initUser(db)
	if err != nil {
		return err
	}
	
	return nil
}