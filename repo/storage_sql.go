package repo

import (
	"context"
	"database/sql"
)

type Storage interface {
	Add(ctx context.Context, qty, id int) error
	List(ctx context.Context) ([]Ingredient, error)
}

type Ingredient struct {
	ID   int    `json:"id"`
	Qty  int    `json:"qty"`
	Name string `json:"name"`
}

type storageSQL struct {
	db *sql.DB
}

func NewStorageSQL(db *sql.DB) *storageSQL {
	return &storageSQL{db: db}
}

// func (s *storageSQL) Add(idx int, qty int) Ingredient {
// 	s.ingredients[idx].Qty += qty

// 	return s.ingredients[idx]
// }

func (s *storageSQL) List(ctx context.Context) ([]Ingredient, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, qty, name FROM storage")
	if err != nil {
		return nil, err
	}

	var stList []Ingredient
	var ingredient Ingredient
	for rows.Next() {
		err := rows.Scan(&ingredient.ID, &ingredient.Qty, &ingredient.Name)
		if err != nil {
			return nil, err
		}

		stList = append(stList, ingredient)
	}

	return stList, nil
}

func (s *storageSQL) Add(ctx context.Context, qty, id int) error {
	_, err := s.db.ExecContext(ctx, "UPDATE storage SET qty = qty + ? WHERE id = ? ", qty, id)
	if err != nil {
		return err
	}

	return nil
}
