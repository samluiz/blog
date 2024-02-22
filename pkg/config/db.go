package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var createTableStatement = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT 0,
    avatar TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    content TEXT DEFAULT '',
    tags TEXT DEFAULT '',
    author_id INTEGER NOT NULL,
    visibility TEXT DEFAULT 'PRIVATE',
    is_published BOOLEAN DEFAULT FALSE,
    published_at TEXT DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT DEFAULT '',
    post_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

func NewConnection() (*sqlx.DB, error) {

	db, err := sqlx.Connect("sqlite", "data.db")

	if err != nil {
		return nil, err
	}

	err = initTables(db)

	if err != nil {
		return nil, err
	}

	return db, nil
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

	_, err := db.Exec(createTableStatement)
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
