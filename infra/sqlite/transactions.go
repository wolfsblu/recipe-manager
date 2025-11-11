package sqlite

import (
	"context"
	"database/sql"

	"github.com/wolfsblu/recipe-manager/domain"
)

// TxStore wraps a Store to provide transaction-scoped database operations
type TxStore struct {
	*Store
	tx *sql.Tx
}

// DB returns the transaction as a *sql.DB for go-jet compatibility
// Note: go-jet's Query and Exec methods work with both *sql.DB and *sql.Tx
// through their underlying interfaces
func (t *TxStore) DB() interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
} {
	return t.tx
}

// WithTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (s *Store) WithTransaction(ctx context.Context, fn func(*TxStore) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.WrapError(domain.ErrStartingTransaction, err)
	}

	txStore := &TxStore{
		Store: s,
		tx:    tx,
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(txStore); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return domain.WrapError(domain.ErrCommittingTransaction, err)
	}
	return nil
}
