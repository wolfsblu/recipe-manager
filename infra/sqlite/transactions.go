package sqlite

import (
	"context"

	"github.com/wolfsblu/go-chef/domain"
)

func (s *Store) Begin(ctx context.Context) error {
	if s.tx != nil {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.WrapError(domain.ErrStartingTransaction, err)
	}
	s.tx = tx
	s.qtx = s.q.WithTx(tx)
	return nil
}

func (s *Store) Commit() error {
	if s.tx == nil {
		return nil
	}
	defer func() {
		s.qtx = nil
		s.tx = nil
	}()
	err := s.tx.Commit()
	if err != nil {
		return domain.WrapError(domain.ErrCommittingTransaction, err)
	}
	return nil
}

func (s *Store) Rollback() {
	if s.tx == nil {
		return
	}
	defer func() {
		s.qtx = nil
		s.tx = nil
	}()
	_ = s.tx.Rollback()
}

func (s *Store) query() *Queries {
	if s.qtx != nil {
		return s.qtx
	}
	return s.q
}
