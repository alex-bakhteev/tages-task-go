package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"tages-task-go/internal/usecase"
)

type Handler struct {
	storeUC usecase.StoreUseCase
}

func NewHandler(storeUC usecase.StoreUseCase) *Handler {
	return &Handler{
		storeUC: storeUC,
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
