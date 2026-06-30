package products

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"marketplace/backend/internal/auth"
	"marketplace/backend/internal/httpapi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type productResponse struct {
	ID          string `json:"id"`
	SellerID    string `json:"sellerId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	sellerID, ok := auth.SellerIDFromContext(r.Context())
	if !ok {
		httpapi.WriteError(w, http.StatusUnauthorized, "unauthenticated", "Authentication required", nil)
		return
	}
	var body createProductRequest
	if err := httpapi.DecodeJSON(r, &body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON", nil)
		return
	}
	product, err := h.service.CreateDraft(r.Context(), sellerID, CreateProductInput{
		Name:        body.Name,
		Description: body.Description,
		Status:      body.Status,
	})
	if err != nil {
		writeProductError(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusCreated, toProductResponse(product))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	sellerID, ok := auth.SellerIDFromContext(r.Context())
	if !ok {
		httpapi.WriteError(w, http.StatusUnauthorized, "unauthenticated", "Authentication required", nil)
		return
	}
	products, err := h.service.ListSellerProducts(r.Context(), sellerID)
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
		return
	}
	response := make([]productResponse, 0, len(products))
	for _, product := range products {
		response = append(response, toProductResponse(product))
	}
	httpapi.WriteJSON(w, http.StatusOK, map[string][]productResponse{"products": response})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	sellerID, ok := auth.SellerIDFromContext(r.Context())
	if !ok {
		httpapi.WriteError(w, http.StatusUnauthorized, "unauthenticated", "Authentication required", nil)
		return
	}
	productID := chi.URLParam(r, "productID")
	product, err := h.service.GetSellerDraft(r.Context(), sellerID, productID)
	if err != nil {
		writeProductError(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, toProductResponse(product))
}

func writeProductError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidProductName):
		httpapi.WriteError(w, http.StatusBadRequest, "validation_error", "Product name must contain between 3 and 200 characters", map[string]string{
			"name": "must contain between 3 and 200 characters",
		})
	case errors.Is(err, ErrProductNotFound):
		httpapi.WriteError(w, http.StatusNotFound, "not_found", "Product not found", nil)
	default:
		httpapi.WriteError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
	}
}

func toProductResponse(product Product) productResponse {
	return productResponse{
		ID:          product.ID,
		SellerID:    product.SellerID,
		Name:        product.Name,
		Description: product.Description,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   product.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	}
}
