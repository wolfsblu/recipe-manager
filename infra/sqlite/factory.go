package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/wolfsblu/go-chef/infra/env"
	"sync"
)

type Store struct {
	db   *sql.DB
	path string
	q    *Queries
	qtx  *Queries
	tx   *sql.Tx
}

var (
	store *Store
	once  sync.Once
)

func NewSqliteStore() (*Store, error) {
	var err error
	dbPath := env.MustGet("DB_PATH")
	once.Do(func() {
		con, err := connect(dbPath)
		if err == nil {
			store = &Store{db: con, q: New(con), path: dbPath}
			err = store.migrate()
		}
	})
	return store, err
}

func connect(path string) (*sql.DB, error) {
	constr := fmt.Sprintf("%s?_pragma=foreign_keys(1)", path)
	con, err := sql.Open("sqlite", constr)
	if err != nil {
		return nil, err
	}
	return con, nil
}
