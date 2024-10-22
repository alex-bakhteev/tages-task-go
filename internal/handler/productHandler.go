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

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// Хендлер для получения всех продуктов
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAllProducts(r.Context())
	if err != nil {
		logger.Log.Println("Error fetching products:", err) // Use logger
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":  "Products retrieved successfully",
		"products": products,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Хендлер для получения продукта по ID
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "Product retrieved successfully",
		"product": product,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Хендлер для создания нового продукта
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.repo.CreateProduct(r.Context(), product)
	if err != nil {
		logger.Log.Println("Error creating product:", err) // Use logger
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Product created successfully",
		"product": product,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
