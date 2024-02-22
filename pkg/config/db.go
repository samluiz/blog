package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/libsql/go-libsql"
	"golang.org/x/crypto/bcrypt"
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
	url := os.Getenv("DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	localDbName := "local.db"

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, localDbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, url, authToken)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()

	db := sql.OpenDB(connector)
	defer db.Close()

	dbWrapper := sqlx.NewDb(db, "libsql")

	err = initTables(dbWrapper)

	if err != nil {
		return nil, err
	}

	return dbWrapper, nil
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
	// schema, err := loadSchema()
	// if err != nil {
	// 	log.Default().Printf("Error loading schema: %v", err)
	// 	return err
	// }

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
