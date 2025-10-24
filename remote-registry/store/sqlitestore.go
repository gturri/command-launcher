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
	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("SQLite store initialized at", dbPath)
	store := &SqliteStore{
		db: db,
	}
	store.createTablesIfNotExists()
	return store, nil
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

func (s *SqliteStore) createTablesIfNotExists() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS registries (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			description TEXT
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS registry_admins (
			registry_id INTEGER,
			admin TEXT,
			FOREIGN KEY (registry_id) REFERENCES registries(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS registry_custom_values (
			registry_id INTEGER,
			key TEXT,
			value TEXT,
			FOREIGN KEY (registry_id) REFERENCES registries(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS packages (
			id INTEGER PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			description TEXT,
			registry_id INTEGER,
			FOREIGN KEY (registry_id) REFERENCES registries(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS package_custom_values (
			package_id INTEGER,
			key TEXT,
			value TEXT,
			FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS package_versions (
			package_id INTEGER,
			version TEXT,
			url TEXT,
			checksum TEXT,
			start_partition INT,
			end_partition INT,
			FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS package_admins (
			package_id INTEGER,
			admin TEXT,
			FOREIGN KEY (package_id) REFERENCES packages(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS deleted_packages (
			registry_id INTEGER,
			package_name TEXT,
			FOREIGN KEY (registry_id) REFERENCES registries(id) ON DELETE CASCADE
			PRIMARY KEY (registry_id, package_name)
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS deleted_registries (
			registry_name TEXT,
			PRIMARY KEY (registry_name)
		);
	`)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *SqliteStore) Close() error {
	return s.db.Close()
}
