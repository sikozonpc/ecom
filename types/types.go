package types

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Address   string    `json:"address"` //workinprogress
	CreatedAt time.Time `json:"createdAt"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	// note that this isn't the best way to handle quantity
	// because it's not atomic (in ACID), but it's good enough for this example
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

type CartCheckoutItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

// aggregate
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

// contains the order history
type OrderHistory struct {
	Orders []Order `json:"orders"`
}
type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	GetOrderItemByOrderID(id int) ([]OrderItem, error)
	GetOrdersByUserID(userID int) ([]Order, error)
}

type ProductStore interface {
	GetProductByID(id int) (*Product, error)
	GetProductsByID(ids []int) ([]Product, error)
	GetProducts() ([]*Product, error)
	CreateProduct(CreateProductPayload) error
	UpdateProduct(Product) error
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
	GetOrdersByID(orderID int) (Order, error)
	UpdateOrderStatus(orderID int, status string) error
}
type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
	Address   string `json:"address" validate:"required"` //address on signup where product will be deliverd to
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CartCheckoutPayload struct {
	Items []CartCheckoutItem `json:"items" validate:"required"`
	//address to be provided on check out or if the user
	//wan to use the same address as they signed up with
	Address string `json:"address"`
}
