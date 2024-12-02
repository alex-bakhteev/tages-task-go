package http

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"tages-task-go/internal/models/modelsuc"
	"tages-task-go/internal/usecases"
	"tages-task-go/pkg/logging"
)

//go:generate go run github.com/vektra/mockery/v2@v2.49.1 --name=ProductUsecase
type ProductUsecase interface {
	CreateProduct(ctx context.Context, product modelsuc.ProductUC) error
	GetProduct(ctx context.Context, id int) (modelsuc.ProductUC, error)
	GetAllProducts(ctx context.Context) ([]modelsuc.ProductUC, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.49.1 --name=OrderUsecase
type OrderUsecase interface {
	CreateOrder(ctx context.Context, order modelsuc.OrderUC) error
	GetOrder(ctx context.Context, id int) (modelsuc.OrderUC, error)
	GetAllOrders(ctx context.Context) ([]modelsuc.OrderUC, error)
}

type Handler struct {
	Usecase usecases.Usecase
	Logger  *logging.Logger
}

func NewHandler(usecase *usecases.Usecase, logger *logging.Logger) *Handler {
	return &Handler{
		Usecase: *usecase,
		Logger:  logger,
	}
}

// InitRoutes инициализирует маршруты для всех сущностей
func (h *Handler) InitRoutes(logger *logging.Logger) *mux.Router {
	router := mux.NewRouter()

	// Middleware с логгером
	router.Use(WithLoggerMiddleware(logger))

	// Роуты
	h.registerProductRoutes(router)
	h.registerOrderRoutes(router)

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

// Middleware для добавления x-request-id
func WithLoggerMiddleware(logger *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлекаем x-request-id или генерируем новый
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Добавляем x-request-id в логгер
			requestLogger := logger.GetLoggerWithField("x-request-id", requestID)

			// Добавляем логгер и x-request-id в контекст
			ctx := logging.WithRequestID(context.Background(), requestID)
			ctx = logging.AddLoggerToContext(ctx, requestLogger)

			// Обновляем запрос с новым контекстом
			r = r.WithContext(ctx)

			// Логируем начало запроса
			requestLogger.Infof("Starting request: %s %s", r.Method, r.URL.Path)

			next.ServeHTTP(w, r)
		})
	}
}
