package postgresql

import (
	"context"
	"fmt"
	"log"
	"tages-task-go/pkg/logging"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"tages-task-go/internal/config"
)

var dbPool *pgxpool.Pool

func InitDB(logger *logging.Logger) *pgxpool.Pool {
	cfg := config.GetConfig()
	logger.Info("Connecting to PostgreSQL...")

	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.Storage.Username,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Database,
	)

	// Формируем строку подключения для миграций
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Storage.Username,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Database,
	)

	// Создаем конфигурацию пула
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Failed to parse database configuration: %v", err)
	}

	// Настраиваем пул соединений
	poolConfig.MaxConns = 10
	poolConfig.HealthCheckPeriod = time.Minute

	// Инициализируем пул соединений
	dbPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Failed to create a connection pool: %v", err)
	}

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	// Запуск миграций
	RunMigrations(databaseURL)

	logger.Info("Connected to PostgreSQL")
	return dbPool
}
