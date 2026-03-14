package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"titanbay/internal/domain"
	"titanbay/internal/service"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service *service.App
}

func NewHandler(service *service.App) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// /funds
	mux.HandleFunc("GET /funds", h.handleGetFunds)
	mux.HandleFunc("POST /funds", h.handleCreateFund)
	mux.HandleFunc("PUT /funds", h.handleUpdateFund)
	mux.HandleFunc("GET /funds/{id}", h.handleGetFundByID)

	// /investors
	mux.HandleFunc("GET /investors", h.handleGetInvestors)
	mux.HandleFunc("POST /investors", h.handleCreateInvestor)

	// /investments
	mux.HandleFunc("GET /funds/{fund_id}/investments", h.handleGetFundInvestments)
	mux.HandleFunc("POST /funds/{fund_id}/investments", h.handleCreateInvestment)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

// List all funds
func (h *Handler) handleGetFunds(w http.ResponseWriter, r *http.Request) {
	funds, err := h.service.Funds.GetAllFunds()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, funds)
}

// Create a new fund
func (h *Handler) handleCreateFund(w http.ResponseWriter, r *http.Request) {
	var input domain.Fund

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newFund, err := h.service.Funds.CreateFund(input)
	if err != nil {
		http.Error(w, "Could not create fund", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, newFund)
}

// Update existing fund
func (h *Handler) handleUpdateFund(w http.ResponseWriter, r *http.Request) {
	var input domain.Fund

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if input.ID == uuid.Nil {
		http.Error(w, "Missing 'id' in request body", http.StatusBadRequest)
		return
	}

	updatedFund, err := h.service.Funds.UpdateFund(input.ID, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Fund not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, updatedFund)
}

// Get specific fund
func (h *Handler) handleGetFundByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	fund, err := h.service.Funds.GetFundByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Fund not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, fund)
}

// List all investors
func (h *Handler) handleGetInvestors(w http.ResponseWriter, r *http.Request) {
	investors, err := h.service.Investors.GetAllInvestors()
	if err != nil {
		http.Error(w, "Failed to retrieve investors", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, investors)
}

// Create a new investor
func (h *Handler) handleCreateInvestor(w http.ResponseWriter, r *http.Request) {
	var input domain.Investor

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Name == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	newInvestor, err := h.service.Investors.CreateInvestor(input)
	if err != nil {
		http.Error(w, "Could not create investor", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, newInvestor)
}

// List all investments for a specific fund
func (h *Handler) handleGetFundInvestments(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("fund_id")
	fundID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid fund UUID", http.StatusBadRequest)
		return
	}

	investments, err := h.service.Investments.GetInvestmentsByFund(fundID)
	if err != nil {
		http.Error(w, "Failed to retrieve investments", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, investments)
}

// Create a new investment to fund
func (h *Handler) handleCreateInvestment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("fund_id")
	fundID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid fund UUID", http.StatusBadRequest)
		return
	}

	var input domain.Investment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input.FundID = fundID

	newInv, err := h.service.Investments.CreateInvestment(input)
	if err != nil {
		http.Error(w, "Could not record investment", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, newInv)
}
