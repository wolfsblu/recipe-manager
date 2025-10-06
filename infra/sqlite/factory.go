package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/wolfsblu/recipe-manager/infra/env"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/mapper"
)

type Store struct {
	db     *sql.DB
	path   string
	q      *database.Queries
	mapper *mapper.DBMapper
}

// NewSqliteStore creates a new SQLite store instance with proper connection pooling.
// The caller is responsible for calling Close() when done with the store.
func NewSqliteStore() (*Store, error) {
	dbPath := env.MustGet("DB_PATH")
	con, err := connect(dbPath)
	if err != nil {
		return nil, err
	}
	con.SetMaxOpenConns(1)

	store := &Store{
		db:     con,
		mapper: mapper.New(),
		path:   dbPath,
		q:      database.New(con),
	}

	if err := store.migrate(); err != nil {
		_ = con.Close()
		return nil, err
	}

	return store, nil
}

// Close closes the database connection and releases all resources.
// It should be called when the store is no longer needed.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func connect(path string) (*sql.DB, error) {
	constr := fmt.Sprintf("%s?_fk=1&_journal=WAL&_sync=NORMAL", path)
	con, err := sql.Open("sqlite", constr)
	if err != nil {
		return nil, err
	}
	return con, nil
}
