package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sikozonpc/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProductByID(productID int) (*types.Product, error) {
	rows, err := s.db.Query("SELECT p.id, p.name, p.description, p.image, p.price, pq.quantity, p.createdat FROM products p join productquantity pq on p.id = pq.id WHERE p.id = ?", productID)
	if err != nil {
		return nil, err
	}

	p := new(types.Product)
	for rows.Next() {
		p, err = scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (s *Store) GetProductsByID(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT p.id, p.name, p.description, p.image, p.price, pq.quantity, p.createdat FROM products p join productquantity pq on p.id = pq.id WHERE p.id IN (?%s)", placeholders)

	// Convert productIDs to []interface{}
	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil

}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT p.id, p.name, p.description, p.image, p.price, pq.quantity, p.createdat FROM products p join productquantity pq on p.id = pq.id")
	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {
	
	// transaction to add product and quantity together
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT INTO products (name, price, image, description) VALUES (?, ?, ?, ?)", product.Name, product.Price, product.Image, product.Description)
	if err != nil {
		tx.Rollback()
		return err
	}

	productID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO productquantity (id, quantity) VALUES (?, ?)", productID, product.Quantity)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateProduct(product types.Product) error {

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}


	_, err = tx.Exec("UPDATE products SET name = ?, price = ?, image = ?, description = ? WHERE id = ?", product.Name, product.Price, product.Image, product.Description, product.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE productquantity SET quantity = ? WHERE id = ?", product.Quantity, product.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateProductQuantity(product types.Product) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE productquantity SET quantity = ? WHERE id = ?", product.Quantity, product.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}
