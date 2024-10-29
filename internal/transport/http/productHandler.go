package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tages-task-go/pkg/models"
	"tages-task-go/pkg/models/transport"
)

func (h *Handler) registerProductRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.createProduct).Methods("POST")
	router.HandleFunc("/products", h.getAllProducts).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", h.getProductByID).Methods("GET")
}

// createProduct - обработчик для создания нового продукта
func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var productDTO transport.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&productDTO); err != nil {
		handleError(w, err, "Invalid request payload", http.StatusBadRequest)
		return
	}

	productUC := models.FromDtoToUseCaseProduct(productDTO)
	if err := h.storeUC.CreateProduct(r.Context(), productUC); err != nil {
		handleError(w, err, "Failed to create product", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{"message": "Product created successfully"})
}

// getAllProducts - обработчик для получения всех продуктов
func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	productsUC, err := h.storeUC.GetAllProducts(r.Context())
	if err != nil {
		handleError(w, err, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	var productsDTO []transport.ProductDTO
	for _, productUC := range productsUC {
		productsDTO = append(productsDTO, models.FromUseCaseToDtoProduct(productUC))
	}

	sendJSONResponse(w, http.StatusOK, productsDTO)
}

// getProduct - обработчик для получения продукта по ID
func (h *Handler) getProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(w, err, "Invalid product ID", http.StatusBadRequest)
		return
	}

	productUC, err := h.storeUC.GetProduct(r.Context(), id)
	if err != nil {
		handleError(w, err, "Product not found", http.StatusNotFound)
		return
	}

	productDTO := models.FromUseCaseToDtoProduct(productUC)
	sendJSONResponse(w, http.StatusOK, productDTO)
}
