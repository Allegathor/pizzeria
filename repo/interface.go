package repo

import "context"

type Storage interface {
	Add(ctx context.Context, qty, id int) error
	List(ctx context.Context) ([]Ingredient, error)
}
