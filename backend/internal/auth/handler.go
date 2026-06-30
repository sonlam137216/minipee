package auth

import (
	"errors"
	"net/http"

	"marketplace/backend/internal/httpapi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type authRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type sellerResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type authResponse struct {
	Seller      sellerResponse `json:"seller"`
	AccessToken string         `json:"accessToken"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body authRequest
	if err := httpapi.DecodeJSON(r, &body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON", nil)
		return
	}
	result, err := h.service.Register(r.Context(), RegisterInput{
		Email:       body.Email,
		Password:    body.Password,
		DisplayName: body.DisplayName,
	})
	if err != nil {
		writeAuthError(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusCreated, toAuthResponse(result))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body authRequest
	if err := httpapi.DecodeJSON(r, &body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, "invalid_json", "Request body must be valid JSON", nil)
		return
	}
	result, err := h.service.Login(r.Context(), LoginInput{Email: body.Email, Password: body.Password})
	if err != nil {
		writeAuthError(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, toAuthResponse(result))
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidAuthInput):
		httpapi.WriteError(w, http.StatusBadRequest, "validation_error", "Invalid registration data", nil)
	case errors.Is(err, ErrEmailAlreadyExists):
		httpapi.WriteError(w, http.StatusConflict, "email_already_exists", "Email is already registered", nil)
	case errors.Is(err, ErrInvalidCredentials):
		httpapi.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password", nil)
	default:
		httpapi.WriteError(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
	}
}

func toAuthResponse(result AuthResult) authResponse {
	return authResponse{Seller: toSellerResponse(result.Seller), AccessToken: result.AccessToken}
}

func toSellerResponse(seller Seller) sellerResponse {
	return sellerResponse{
		ID:          seller.ID,
		Email:       seller.Email,
		DisplayName: seller.DisplayName,
		CreatedAt:   seller.CreatedAt.UTC().Format(timeFormat),
		UpdatedAt:   seller.UpdatedAt.UTC().Format(timeFormat),
	}
}

const timeFormat = "2006-01-02T15:04:05Z07:00"
