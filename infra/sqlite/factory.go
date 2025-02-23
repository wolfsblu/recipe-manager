package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/google/wire"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/env"
	"log"
)

type Store struct {
	db   *sql.DB
	path string
	q    *Queries
	qtx  *Queries
	tx   *sql.Tx
}

var Set = wire.NewSet(
	NewSqliteStore,
	wire.Bind(new(domain.RecipeStore), new(*Store)),
)

func NewSqliteStore() (*Store, error) {
	log.Println("creating new sqlite store")
	dbPath := env.MustGet("DB_PATH")
	con, err := connect(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	s := &Store{db: con, q: New(con), path: dbPath}
	err = s.migrate()
	if err != nil {
		log.Fatalln("failed to apply database migrations:", err)
	}

	return s, nil
}

func connect(path string) (*sql.DB, error) {
	constr := fmt.Sprintf("%s?_pragma=foreign_keys(1)", path)
	con, err := sql.Open("sqlite", constr)
	if err != nil {
		return nil, err
	}
	return con, nil
}
