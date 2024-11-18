package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"tages-task-go/internal/models/modelsuc"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, product modelsuc.ProductUC) error
	GetProduct(ctx context.Context, id int) (modelsuc.ProductUC, error)
	GetAllProducts(ctx context.Context) ([]modelsuc.ProductUC, error)
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order modelsuc.OrderUC) error
	GetOrder(ctx context.Context, id int) (modelsuc.OrderUC, error)
	GetAllOrders(ctx context.Context) ([]modelsuc.OrderUC, error)
}

type Handler struct {
	ProductUsecase ProductUsecase
	OrderUsecase   OrderUsecase
}

func NewHandler(productUC ProductUsecase, orderUC OrderUsecase) *Handler {
	return &Handler{
		ProductUsecase: productUC,
		OrderUsecase:   orderUC,
	}
}

// InitRoutes инициализирует маршруты для всех сущностей
func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Подключаем маршруты для Order
	h.registerOrderRoutes(router)

	// Подключаем маршруты для Product
	h.registerProductRoutes(router)

	return router
}

// sendJSONResponse отправляет JSON-ответ с указанным статусом
func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// handleError обрабатывает ошибки и отправляет соответствующий HTTP-ответ
func handleError(w http.ResponseWriter, err error, msg string, status int) {
	if err != nil {
		http.Error(w, msg, status)
	}
}
