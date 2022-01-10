package repo

import (
	"context"
	"database/sql"
)

const UpdateQtyQuery string = `
	UPDATE storage s
	INNER JOIN recipe_ingredients r
	ON s.id = r.storage_id
	SET qty = qty - req_qty
	WHERE recipe_id = (
		SELECT id
		FROM recipes
		WHERE id = (
				SELECT recipe_id
				FROM menu WHERE id = ?)) AND req_qty < qty;
`

type Order interface {
	Create(ctx context.Context, order OrderInfo) error
	List(ctx context.Context) ([]OrderDetails, error)
}

type OrderInfo struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Street    string `json:"street"`
	House     string `json:"house"`
	Housing   string `json:"housing"`
	Apartment string `json:"apartment"`
	PizzaID   int    `json:"pizza_id"`
}

type OrderDetails struct {
	ID     int `json:"id"`
	Status int `json:"status"`
	OrderInfo
}

type orderSQL struct {
	db *sql.DB
}

func NewOrderSQL(db *sql.DB) *orderSQL {
	return &orderSQL{db: db}
}

func (o *orderSQL) Create(ctx context.Context, info OrderInfo) error {
	result, err := o.db.ExecContext(ctx, UpdateQtyQuery, info.PizzaID)
	if err != nil {
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowCount == 0 {
		return &StorageError{}
	}

	_, err = o.db.ExecContext(
		ctx,
		`INSERT INTO orders SET 
		name = ?, phone = ?, street = ?, house = ?, housing = ?, apartment = ?, pizza_id = ?`,
		info.Name, info.Phone, info.Street, info.House, info.Housing, info.Apartment, info.PizzaID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (o *orderSQL) List(ctx context.Context) ([]OrderDetails, error) {
	rows, err := o.db.QueryContext(ctx, "SELECT id, name, phone, street, house, housing, apartment, pizza_id, status FROM orders")
	if err != nil {
		return nil, err
	}

	var orders []OrderDetails
	var order OrderDetails
	for rows.Next() {
		err := rows.Scan(&order.ID, &order.Name, &order.Phone, &order.Street, &order.House, &order.Housing, &order.Apartment, &order.PizzaID, &order.Status)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

type StorageError struct{}

func (m *StorageError) Error() string {
	return "Wrong pizzaID or you need to resupply storage"
}
