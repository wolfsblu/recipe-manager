package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (s *Store) GetUnits(ctx context.Context, req domain.Page) (domain.Result[domain.Unit], error) {
	cursor, err := domain.DecodeCursor[*domain.NameCursor](req.Cursor)
	if err != nil {
		cursor = &domain.NameCursor{}
	}
	result, err := s.query().GetUnits(ctx, database.GetUnitsParams{
		LastName: cursor.LastName,
		LastID:   cursor.LastID,
		Limit:    int64(req.Limit + 1),
	})
	if err != nil {
		return domain.Result[domain.Unit]{}, err
	}

	units := make([]domain.Unit, 0, len(result))
	for _, unit := range result {
		units = append(units, s.mapper.ToUnit(unit))
	}

	return domain.NewPagedResult(units, req.Limit, func(u domain.Unit) domain.NameCursor {
		return domain.NameCursor{
			LastID:   u.ID,
			LastName: u.Name,
		}
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
