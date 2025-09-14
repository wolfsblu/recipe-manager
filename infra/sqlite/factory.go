package sqlite

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/wolfsblu/recipe-manager/infra/env"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/mapper"
)

type Store struct {
	db     *sql.DB
	path   string
	q      *database.Queries
	qtx    *database.Queries
	tx     *sql.Tx
	mapper *mapper.DBMapper
}

var (
	store *Store
	once  sync.Once
)

func NewSqliteStore() (*Store, error) {
	var err error
	dbPath := env.MustGet("DB_PATH")
	once.Do(func() {
		var con *sql.DB
		con, err = connect(dbPath)
		if err == nil {
			store = &Store{
				db:     con,
				mapper: mapper.New(),
				path:   dbPath,
				q:      database.New(con),
			}
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
