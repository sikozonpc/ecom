package txManager

import (
	"context"
	"database/sql"

	"github.com/sikozonpc/ecom/services/order"
	"github.com/sikozonpc/ecom/services/product"
	"github.com/sikozonpc/ecom/services/user"
	"github.com/sikozonpc/ecom/types"
)

type TxContext interface {
	UserStore() types.UserStore
	ProductStore() types.ProductStore
	OrderStore() types.OrderStore
}

type sqlTxContext struct {
	userStore    types.UserStore
	productStore types.ProductStore
	orderStore   types.OrderStore
}

func (c *sqlTxContext) UserStore() types.UserStore {
	return c.userStore
}
func (c *sqlTxContext) ProductStore() types.ProductStore {
	return c.productStore
}
func (c *sqlTxContext) OrderStore() types.OrderStore {
	return c.orderStore
}

type TransactionManager interface {
	WithCtx(ctx context.Context, fn func(TxContext) error) error
}
type TrxManager struct {
	db *sql.DB
}

func NewTrxManager(db *sql.DB) *TrxManager {
	return &TrxManager{db: db}
}

func (s *TrxManager) WithCtx(ctx context.Context, fn func(TxContext) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txCx := &sqlTxContext{
		userStore:    user.NewStore(tx),
		productStore: product.NewStore(tx),
		orderStore:   order.NewStore(tx),
	}
	if err := fn(txCx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
