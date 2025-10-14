package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (s *Store) GetUnits(ctx context.Context) ([]domain.Unit, error) {
	result, err := s.query().GetUnits(ctx)
	if err != nil {
		return nil, err
	}

	units := make([]domain.Unit, len(result))
	for i, unit := range result {
		units[i] = s.mapper.ToUnit(unit)
	}

	return units, nil
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
