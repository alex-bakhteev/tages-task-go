package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tages-task-go/internal/logger" // Import your logger package
	"tages-task-go/internal/model"
	"tages-task-go/internal/repository"
)

type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// Хендлер для получения всех заказов
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.GetOrders(r.Context())
	if err != nil {
		logger.Log.Println("Error fetching orders:", err) // Use logger
		http.Error(w, "Error fetching orders", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Orders retrieved successfully",
		"orders":  orders,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Хендлер для получения заказа по ID
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.repo.GetOrderByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "Order retrieved successfully",
		"order":   order,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Хендлер для создания нового заказа
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.repo.CreateOrder(r.Context(), order)
	if err != nil {
		logger.Log.Println("Error creating order:", err) // Use logger
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Order created successfully",
		"order":   order,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
