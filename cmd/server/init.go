package server

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/internal/service/db/postgresql"
	"tages-task-go/internal/transport/http"
	"tages-task-go/internal/usecase"
	"tages-task-go/pkg/logging"
)

var DbPool *pgxpool.Pool

// Initialize инициализирует необходимые компоненты приложения
func Initialize() (*mux.Router, error) {
	logger := logging.GetLogger()
	logger.Infof("Initializing server")
	// Подключение к базе данных
	DbPool = postgresql.InitDB(logger)

	// Инициализация репозиториев и юзкейсов
	productRepo := postgresql.NewProductRepository(DbPool, logger)
	orderRepo := postgresql.NewOrderRepository(DbPool, logger)
	productUC := usecase.NewProductUseCase(productRepo, logger)
	orderUC := usecase.NewOrderUseCase(orderRepo, logger)
	storeUC := http.NewStoreUseCase(orderUC, productUC)

	// Инициализация хендлеров и маршрутов
	handler := http.NewHandler(storeUC)
	router := handler.InitRoutes()

	return router, nil
}
