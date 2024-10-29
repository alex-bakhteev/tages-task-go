package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tages-task-go/pkg/models"
	"tages-task-go/pkg/models/transport"
)

func (h *Handler) registerOrderRoutes(router *mux.Router) {
	router.HandleFunc("/orders", h.createOrder).Methods("POST")
	router.HandleFunc("/orders", h.getAllOrders).Methods("GET")
	router.HandleFunc("/orders/{id:[0-9]+}", h.getOrderByID).Methods("GET")
}

// createOrder - обработчик для создания нового заказа
func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var orderDTO transport.OrderDTO
	if err := json.NewDecoder(r.Body).Decode(&orderDTO); err != nil {
		handleError(w, err, "Invalid request payload", http.StatusBadRequest)
		return
	}

	orderUC := models.FromDtoToUseCaseOrder(orderDTO)
	if err := h.storeUC.CreateOrder(r.Context(), orderUC); err != nil {
		handleError(w, err, "Failed to create order", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{"message": "Order created successfully"})
}

// getOrders - обработчик для получения всех заказов
func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	ordersUC, err := h.storeUC.GetAllOrders(r.Context())
	if err != nil {
		handleError(w, err, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	var ordersDTO []transport.OrderDTO
	for _, orderUC := range ordersUC {
		ordersDTO = append(ordersDTO, models.FromUseCaseToDtoOrder(orderUC))
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

	orderUC, err := h.storeUC.GetOrder(r.Context(), id)
	if err != nil {
		handleError(w, err, "Order not found", http.StatusNotFound)
		return
	}

	orderDTO := models.FromUseCaseToDtoOrder(orderUC)
	sendJSONResponse(w, http.StatusOK, orderDTO)
}
