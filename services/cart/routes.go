package cart

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sikozonpc/ecom/services/auth"
	txManager "github.com/sikozonpc/ecom/services/txmanager"
	"github.com/sikozonpc/ecom/types"
	"github.com/sikozonpc/ecom/utils"
)

type Handler struct {
	store      types.ProductStore
	orderStore types.OrderStore
	userStore  types.UserStore
	txManager  txManager.TransactionManager
}

func NewHandler(
	store types.ProductStore,
	orderStore types.OrderStore,
	userStore types.UserStore,
	txManager txManager.TransactionManager,
) *Handler {
	return &Handler{
		store:      store,
		orderStore: orderStore,
		userStore:  userStore,
		txManager:  txManager,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/order/cancel/{order_id}", auth.WithJWTAuth(h.cancelOrderHandler, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	productIds, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// At this point, there are two possible error sources:
	// 1. Product out of stock — quantity requested exceeds available stock.
	// 2. Database error.
	// For this simple project, returning StatusBadRequest is sufficient.
	// In a more robust project, we'd distinguish between these and handle internal errors separately.
	orderID, totalPrice, err := h.CreateOrder(userID, productIds, cart)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}

func (h *Handler) cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["order_id"]
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order ID"))
		return
	}

	// Get user ID from context (assuming you store it during auth)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	// Call service logic
	if err := h.CancelOrder(userID, orderID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return

	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "order cancelled"})
}
