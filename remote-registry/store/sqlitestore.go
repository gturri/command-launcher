package store

import (
	"database/sql"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(dbPath string) (*SqliteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("SQLite store initialized at", dbPath)
	return &SqliteStore{
		db: db,
	}, nil
}

func (s *SqliteStore) GetSqliteVersion() (string, error) {
	var version string
	err := s.db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("SQLite version:", version)
	return version, nil
}

func (s *SqliteStore) Close() error {
	return s.db.Close()
}
