package product

import (
	"database/sql"

	"github.com/sikozonpc/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProductByID(productID int) (*types.Product, error) {
	return nil, nil
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	return nil, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	return nil
}