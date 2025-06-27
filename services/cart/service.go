package cart

import (
	"context"
	"fmt"

	txManager "github.com/sikozonpc/ecom/services/txmanager"
	"github.com/sikozonpc/ecom/types"
)

func getCartItemsIDs(items []types.CartCheckoutItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productIds[i] = item.ProductID
	}

	return productIds, nil
}

func checkIfCartIsInStock(cartItems []types.CartCheckoutItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	// Instead of returning a single error when an unavailable item is found,
	// consider iterating through the entire cart and checking each item’s availability.
	// If any items are unavailable, build a list or map of errors
	// so the user knows exactly which items are not available.
	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartCheckoutItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}

func (h *Handler) CreateOrder(userID int, productIds []int, cart types.CartCheckoutPayload) (int, float64, error) {
	var orderID int
	var totalPrice float64
	err := h.txManager.WithCtx(context.Background(), func(tc txManager.TxContext) error {
		products, err := tc.ProductStore().GetProductsByID(productIds)
		if err != nil {
			return err
		}
		productsMap := make(map[int]types.Product)
		for _, product := range products {
			productsMap[product.ID] = product
		}

		// check if all products are available
		if err := checkIfCartIsInStock(cart.Items, productsMap); err != nil {
			return err
		}

		// calculate total price
		totalPrice := calculateTotalPrice(cart.Items, productsMap)

		// reduce the quantity of products in the store
		for _, item := range cart.Items {
			product := productsMap[item.ProductID]
			product.Quantity -= item.Quantity
			h.store.UpdateProduct(product)
		}

		// create order record
		orderIDLocal, err := tc.OrderStore().CreateOrder(types.Order{
			UserID:  userID,
			Total:   totalPrice,
			Status:  "pending",
			Address: cart.Address, // user submits addres with the cart details for delivery
		})
		if err != nil {
			return err
		}
		orderID = orderIDLocal

		// create order the items records
		for _, item := range cart.Items {
			tc.OrderStore().CreateOrderItem(types.OrderItem{
				OrderID:   orderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     productsMap[item.ProductID].Price,
			})
		}

		return nil
	})
	return orderID, totalPrice, err
}

func (h *Handler) CancelOrder(userID int, orderID int) error {
	return h.txManager.WithCtx(context.Background(), func(tc txManager.TxContext) error {
		// get order
		order, err := tc.OrderStore().GetOrdersByID(orderID)
		if err != nil {
			return fmt.Errorf("order not found: %w", err)
		}

		// Check ownership
		if order.UserID != userID {
			return fmt.Errorf("unauthorized: user does not own this order")
		}

		// Only allow cancel if not shipped
		if order.Status != "pending" && order.Status != "completed" {
			return fmt.Errorf("cannot cancel order that is already shipped or completed")
		}

		// Update order status
		if err := tc.OrderStore().UpdateOrderStatus(orderID, "cancelled"); err != nil {
			return fmt.Errorf("failed to cancel order: %w", err)
		}

		// Get order items
		items, err := tc.UserStore().GetOrderItemByOrderID(orderID)
		if err != nil {
			return fmt.Errorf("failed to fetch order items: %w", err)
		}

		// Restore product quantities
		for _, item := range items {
			product, err := tc.ProductStore().GetProductByID(item.ProductID)
			if err != nil {
				return fmt.Errorf("failed to fetch product: %w", err)
			}
			product.Quantity += item.Quantity
			if err := tc.ProductStore().UpdateProduct(*product); err != nil {
				return fmt.Errorf("failed to update product quantity: %w", err)
			}
		}

		return nil
	})
}
