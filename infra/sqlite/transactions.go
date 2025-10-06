package sqlite

import (
	"context"
	"database/sql"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

// TxStore wraps a Store with isolated transaction state
// This ensures each transaction has its own state and prevents concurrency issues
type TxStore struct {
	*Store
	tx  *sql.Tx
	qtx *database.Queries
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
		qtx:   s.q.WithTx(tx),
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

// query returns the transaction-aware query executor for TxStore
func (t *TxStore) query() *database.Queries {
	return t.qtx
}

// query returns the regular query executor for Store (when not in a transaction)
func (s *Store) query() *database.Queries {
	return s.q
}
