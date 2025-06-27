package order

import (
	"github.com/sikozonpc/ecom/types"
)

type Store struct {
	db types.DBExecutor
}

func NewStore(db types.DBExecutor) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}

func (s *Store) GetOrdersByID(orderID int) (types.Order, error) {
	row := s.db.QueryRow(`SELECT id, userId, total, status, address, createdAt FROM orders WHERE userId = ?;`, orderID)

	var o types.Order
	err := row.Scan(
		&o.ID,
		&o.UserID,
		&o.Total,
		&o.Status,
		&o.Address,
		&o.CreatedAt,
	)
	if err != nil {
		return types.Order{}, err
	}

	return o, err
}

func (s *Store) UpdateOrderStatus(orderID int, status string) error {
	_, err := s.db.Exec(`UPDATE orders SET status = ? WHERE id = ?`, status, orderID)
	return err
}