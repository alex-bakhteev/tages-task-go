package server

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/internal/service/db/postgresql"
	"tages-task-go/internal/transport/http"
	"tages-task-go/internal/usecases"
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
	usecase := usecases.NewUsecase(productRepo, orderRepo, logger)

	// Инициализация хендлеров и маршрутов
	handler := http.NewHandler(usecase, logger)
	router := handler.InitRoutes(logger)

	return router, nil
}
