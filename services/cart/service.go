package cart

import (
	"fmt"

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

func (h *Handler) createOrder(products []types.Product, cartItems []types.CartCheckoutItem, userID int) (int, float64, error) {
	// create a map of products for easier access
	productsMap := make(map[int]types.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	// check if all products are available
	if err := checkIfCartIsInStock(cartItems, productsMap); err != nil {
		return 0, 0, err
	}

	// calculate total price
	totalPrice := calculateTotalPrice(cartItems, productsMap)

	// reduce the quantity of products in the store
	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.store.UpdateProduct(product)
	}

	// create order record
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address", // could fetch address from a user addresses table
	})
	if err != nil {
		return 0, 0, err
	}

	// create order the items records
	for _, item := range cartItems {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}
