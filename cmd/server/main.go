package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"tages-task-go/internal/db"
	"tages-task-go/internal/handler"
	"tages-task-go/internal/logger" // Import your logger package
	"tages-task-go/internal/repository"
)

func main() {
	logger.InitLogger() // Initialize the logger

	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Error loading .env file") // Use logger instead of log
	}

	// Подключение к базе данных
	dbPool, err := db.ConnectDatabase()
	if err != nil {
		logger.Log.Fatalf("Failed to connect to database: %v", err) // Use logger
	}
	defer dbPool.Close()

	// Репозиторий и хендлеры
	productRepo := repository.NewProductRepository(dbPool)
	productHandler := handler.NewProductHandler(productRepo)
	orderRepo := repository.NewOrderRepository(dbPool)
	orderHandler := handler.NewOrderHandler(orderRepo)

	// Настройка роутера
	r := mux.NewRouter()

	r.HandleFunc("/products", productHandler.GetAllProducts).Methods("GET")
	r.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productHandler.GetProductByID).Methods("GET")
	r.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", orderHandler.GetOrderByID).Methods("GET")
	r.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")

	logger.Log.Fatal(http.ListenAndServe(":8080", r)) // Use logger
}
