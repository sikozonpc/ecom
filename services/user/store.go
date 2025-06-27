package user

import (
	"database/sql"
	"fmt"

	"github.com/sikozonpc/ecom/types"
)

type Store struct {
	db types.DBExecutor
}

func NewStore(db types.DBExecutor) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password, address) VALUES (?, ?, ?, ?,?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.Address)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.Address,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetOrdersByUserID(userID int) ([]types.Order, error) {
	rows, err := s.db.Query(`SELECT id, userId total, status, address, createdAt FROM orders WHERE userId = ?;`, userID)
	if err != nil {
		return nil, err
	}
	var orders = []types.Order{}
	for rows.Next() {
		o := types.Order{}
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Total,
			&o.Status,
			&o.Address,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, err
}

func (r *Store) GetOrderItemByOrderID(orderID int) ([]types.OrderItem, error) {
	query := `
		SELECT 
			oi.id,
			oi.orderId,
			oi.productId,
			p.name,
			p.description,
			p.image,
			oi.quantity,
			oi.price
		FROM order_items oi
		JOIN products p ON oi.productId = p.id
		WHERE oi.orderId = ?
	`

	rows, err := r.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.OrderItem

	for rows.Next() {
		var item types.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Name,
			&item.Image,
			&item.Quantity,
			&item.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
