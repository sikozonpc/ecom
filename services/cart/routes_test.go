package cart

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sikozonpc/ecom/types"
)

var mockProducts = []types.Product{
	{ID: 1, Name: "product 1", Price: 10, Quantity: 100},
	{ID: 2, Name: "product 2", Price: 20, Quantity: 200},
	{ID: 3, Name: "product 3", Price: 30, Quantity: 300},
	{ID: 4, Name: "empty stock", Price: 30, Quantity: 0},
	{ID: 5, Name: "almost stock", Price: 30, Quantity: 1},
}

func TestCartServiceHandler(t *testing.T) {
	productStore := &mockProductStore{}
	orderStore := &mockOrderStore{}
	handler := NewHandler(productStore, orderStore, nil)

	t.Run("should fail to checkout if the cart items do not exist", func(t *testing.T) {
		payload := types.CartCheckoutPayload{
			Items: []types.CartCheckoutItem{
				{ProductID: 99, Quantity: 100},
			},
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/cart/checkout", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/cart/checkout", handler.handleCheckout).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to checkout if the cart has negative quantities", func(t *testing.T) {
		payload := types.CartCheckoutPayload{
			Items: []types.CartCheckoutItem{
				{ProductID: 1, Quantity: 0}, // invalid quantity
			},
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/cart/checkout", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/cart/checkout", handler.handleCheckout).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to checkout if there is no stock for an item", func(t *testing.T) {
		payload := types.CartCheckoutPayload{
			Items: []types.CartCheckoutItem{
				{ProductID: 4, Quantity: 2},
			},
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/cart/checkout", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/cart/checkout", handler.handleCheckout).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to checkout if there is not enough stock", func(t *testing.T) {
		payload := types.CartCheckoutPayload{
			Items: []types.CartCheckoutItem{
				{ProductID: 5, Quantity: 2},
			},
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/cart/checkout", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/cart/checkout", handler.handleCheckout).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should checkout and calculate the price correctly", func(t *testing.T) {
		payload := types.CartCheckoutPayload{
			Items: []types.CartCheckoutItem{
				{ProductID: 1, Quantity: 10},
				{ProductID: 2, Quantity: 20},
				{ProductID: 5, Quantity: 1},
			},
		}

		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/cart/checkout", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/cart/checkout", handler.handleCheckout).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["total_price"] != 530.0 {
			t.Errorf("expected total price to be 530, got %f", response["total_price"])
		}
	})
}

type mockProductStore struct{}

func (m *mockProductStore) GetProductByID(productID int) (*types.Product, error) {
	return &types.Product{}, nil
}

func (m *mockProductStore) GetProducts() ([]*types.Product, error) {
	return []*types.Product{}, nil
}

func (m *mockProductStore) CreateProduct(product types.CreateProductPayload) error {
	return nil
}

func (m *mockProductStore) GetProductsByID(ids []int) ([]types.Product, error) {
	return mockProducts, nil
}

func (m *mockProductStore) UpdateProduct(product types.Product) error {
	return nil
}

type mockOrderStore struct{}

func (m *mockOrderStore) CreateOrder(order types.Order) (int, error) {
	return 0, nil
}

func (m *mockOrderStore) CreateOrderItem(orderItem types.OrderItem) error {
	return nil
}