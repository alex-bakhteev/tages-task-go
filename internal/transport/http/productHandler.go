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
		h.Logger.ErrorCtx(r.Context(), "Invalid request payload: %v", err)
		return
	}

	productUC := modelsuc.ProductUC{
		ID:    productDTO.ID,
		Name:  productDTO.Name,
		Price: productDTO.Price,
	}

	if err := h.Usecase.CreateProduct(r.Context(), productUC); err != nil {
		handleError(w, err, "Failed to create product", http.StatusInternalServerError)
		h.Logger.ErrorCtx(r.Context(), "Failed to create product: %v", err)
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{"message": "Product created successfully"})
	h.Logger.InfoCtx(r.Context(), "Product created successfully: ID=%d, Name=%s", productDTO.ID, productDTO.Name)
}

// getAllProducts - обработчик для получения всех продуктов
func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	productsUC, err := h.Usecase.GetAllProducts(r.Context())
	if err != nil {
		handleError(w, err, "Failed to fetch products", http.StatusInternalServerError)
		h.Logger.ErrorCtx(r.Context(), "Failed to fetch products: %v", err)
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
	h.Logger.InfoCtx(r.Context(), "Fetched all products, count=%d", len(productsDTO))
}

// getProductByID - обработчик для получения продукта по ID
func (h *Handler) getProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handleError(w, err, "Invalid product ID", http.StatusBadRequest)
		h.Logger.ErrorCtx(r.Context(), "Invalid product ID: %s", vars["id"])
		return
	}

	productUC, err := h.Usecase.GetProduct(r.Context(), id)
	if err != nil {
		handleError(w, err, "Product not found", http.StatusNotFound)
		h.Logger.DebugCtx(r.Context(), "Product not found: ID=%d", id)
		return
	}

	productDTO := modelstr.ProductDTO{
		ID:    productUC.ID,
		Name:  productUC.Name,
		Price: productUC.Price,
	}

	sendJSONResponse(w, http.StatusOK, productDTO)
	h.Logger.InfoCtx(r.Context(), "Fetched product: ID=%d, Name=%s", productDTO.ID, productDTO.Name)
}
