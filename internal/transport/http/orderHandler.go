package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tages-task-go/internal/models/modelstr"
	"tages-task-go/internal/models/modelsuc"
)

func (h *Handler) registerOrderRoutes(router *mux.Router) {
	router.HandleFunc("/orders", h.createOrder).Methods("POST")
	router.HandleFunc("/orders", h.getAllOrders).Methods("GET")
	router.HandleFunc("/orders/{id:[0-9]+}", h.getOrderByID).Methods("GET")
}

// createOrder - обработчик для создания нового заказа
func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var orderDTO modelstr.OrderDTO
	if err := json.NewDecoder(r.Body).Decode(&orderDTO); err != nil {
		handleError(w, err, "Invalid request payload", http.StatusBadRequest)
		return
	}

	orderUC := modelsuc.OrderUC{
		ID:        orderDTO.ID,
		ProductID: orderDTO.ProductID,
		Quantity:  orderDTO.Quantity,
	}
	if err := h.OrderUsecase.CreateOrder(r.Context(), orderUC); err != nil {
		handleError(w, err, "Failed to create order", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{"message": "Order created successfully"})
}

// getOrders - обработчик для получения всех заказов
func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	ordersUC, err := h.OrderUsecase.GetAllOrders(r.Context())
	if err != nil {
		handleError(w, err, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	var ordersDTO []modelstr.OrderDTO
	for _, orderUC := range ordersUC {
		orderDTO := modelstr.OrderDTO{
			ID:        orderUC.ID,
			ProductID: orderUC.ProductID,
			Quantity:  orderUC.Quantity,
		}
		ordersDTO = append(ordersDTO, orderDTO)
	}

	sendJSONResponse(w, http.StatusOK, ordersDTO)
}

// getOrderByID - обработчик для получения заказа по ID
func (h *Handler) getOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, "Invalid order ID", http.StatusBadRequest)
		return
	}

	orderUC, err := h.OrderUsecase.GetOrder(r.Context(), id)
	if err != nil {
		handleError(w, err, "Order not found", http.StatusNotFound)
		return
	}

	orderDTO := modelstr.OrderDTO{
		ID:        orderUC.ID,
		ProductID: orderUC.ProductID,
		Quantity:  orderUC.Quantity,
	}
	sendJSONResponse(w, http.StatusOK, orderDTO)
}
