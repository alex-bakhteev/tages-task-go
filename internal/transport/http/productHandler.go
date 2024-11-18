package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"tages-task-go/internal/models/modelstr"
	"tages-task-go/internal/models/modelsuc"
)

func (h *Handler) registerProductRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.createProduct).Methods("POST")
	router.HandleFunc("/products", h.getAllProducts).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", h.getProductByID).Methods("GET")
}

// createProduct - обработчик для создания нового продукта
func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var productDTO modelstr.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&productDTO); err != nil {
		handleError(w, err, "Invalid request payload", http.StatusBadRequest)
		return
	}

	productUC := modelsuc.ProductUC{
		ID:    productDTO.ID,
		Name:  productDTO.Name,
		Price: productDTO.Price,
	}
	if err := h.ProductUsecase.CreateProduct(r.Context(), productUC); err != nil {
		handleError(w, err, "Failed to create product", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{"message": "Product created successfully"})
}

// getAllProducts - обработчик для получения всех продуктов
func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	productsUC, err := h.ProductUsecase.GetAllProducts(r.Context())
	if err != nil {
		handleError(w, err, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	var productsDTO []modelstr.ProductDTO
	for _, productUC := range productsUC {
		productDTO := modelstr.ProductDTO{
			ID:    productUC.ID,
			Name:  productUC.Name,
			Price: productUC.Price,
		}
		productsDTO = append(productsDTO, productDTO)
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

	productUC, err := h.ProductUsecase.GetProduct(r.Context(), id)
	if err != nil {
		handleError(w, err, "Product not found", http.StatusNotFound)
		return
	}

	productDTO := modelstr.ProductDTO{
		ID:    productUC.ID,
		Name:  productUC.Name,
		Price: productUC.Price,
	}
	sendJSONResponse(w, http.StatusOK, productDTO)
}
