package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/gen/model"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/queries"
)

func (s *Store) GetUnits(ctx context.Context, req domain.Page) (domain.Result[domain.Unit], error) {
	cursor, err := domain.DecodeCursor[*domain.NameCursor](req.Cursor)
	if err != nil {
		cursor = &domain.NameCursor{}
	}

	var result []model.Unit
	err = queries.SelectUnits(cursor.LastID, cursor.LastName, int64(req.Limit+1)).Query(s.DB(), &result)
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
	var result model.Unit
	err := queries.InsertUnit(unit.Name, unit.Symbol).Query(s.DB(), &result)
	if err != nil {
		return domain.Unit{}, err
	}

	return domain.Unit{
		ID:     result.ID,
		Name:   result.Name,
		Symbol: result.Symbol,
	}, nil
}

func (s *Store) UpdateUnit(ctx context.Context, unit domain.Unit) error {
	_, err := queries.UpdateUnit(unit.ID, unit.Name, unit.Symbol).Exec(s.DB())
	return err
}

func (s *Store) DeleteUnit(ctx context.Context, id int64) error {
	_, err := queries.DeleteUnit(id).Exec(s.DB())
	return err
}
