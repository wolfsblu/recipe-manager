package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/domain/pagination"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (s *Store) GetUnits(ctx context.Context, req pagination.Page) (pagination.Result[domain.Unit], error) {
	result, err := s.query().GetUnits(ctx, database.GetUnitsParams{
		Cursor: pagination.DecodeCursor(req.Cursor),
		Limit:  int64(req.Limit + 1),
	})
	if err != nil {
		return pagination.Result[domain.Unit]{}, err
	}

	units := make([]domain.Unit, 0, len(result))
	for _, unit := range result {
		units = append(units, s.mapper.ToUnit(unit))
	}

	return pagination.NewResult(units, req.Limit, func(u domain.Unit) int64 {
		return u.ID
	}), nil
}

func (s *Store) CreateUnit(ctx context.Context, unit domain.Unit) (domain.Unit, error) {
	id, err := s.query().CreateUnit(ctx, database.CreateUnitParams{
		Name:   unit.Name,
		Symbol: unit.Symbol,
	})
	if err != nil {
		return domain.Unit{}, err
	}
	return domain.Unit{
		ID:     id,
		Name:   unit.Name,
		Symbol: unit.Symbol,
	}, nil
}

func (s *Store) UpdateUnit(ctx context.Context, unit domain.Unit) error {
	return s.query().UpdateUnit(ctx, database.UpdateUnitParams{
		Name:   unit.Name,
		Symbol: unit.Symbol,
		ID:     unit.ID,
	})
}

func (s *Store) DeleteUnit(ctx context.Context, id int64) error {
	return s.query().DeleteUnit(ctx, id)
}
