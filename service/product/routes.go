package product

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/zinx110/golang-backend-rest/types"
	"github.com/zinx110/golang-backend-rest/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleCreateProducts).Methods(http.MethodPost)
	router.HandleFunc("/products", h.handleGetAllProducts).Methods(http.MethodGet)
}

func (h *Handler) handleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetAllProducts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve products: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)

}

func (h *Handler) handleCreateProducts(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}
	product, err := h.store.CreateProduct(&payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, product)

}
