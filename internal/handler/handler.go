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

	if input.Name == nil || input.VintageYear == nil || input.TargetSizeUSD == nil || input.Status == nil {
		http.Error(w, "name, vintage_year, target_size_usd and status are required", http.StatusBadRequest)
		return
	}

	if *input.Name == "" || *input.VintageYear < 1900 || *input.TargetSizeUSD <= 0 || *input.Status == "" {
		http.Error(w, "name, vintage_year, target_size_usd and status cannot be empty", http.StatusBadRequest)
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
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.ID == uuid.Nil {
		http.Error(w, "missing 'id' in request body", http.StatusBadRequest)
		return
	}

	updatedFund, err := h.service.Funds.UpdateFund(input.ID, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "fund not found", http.StatusNotFound)
			return
		}
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, updatedFund)
}

// Get specific fund
func (h *Handler) handleGetFundByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	fund, err := h.service.Funds.GetFundByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "fund not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, fund)
}

// List all investors
func (h *Handler) handleGetInvestors(w http.ResponseWriter, r *http.Request) {
	investors, err := h.service.Investors.GetAllInvestors()
	if err != nil {
		http.Error(w, "failed to retrieve investors", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, investors)
}

// Create a new investor
func (h *Handler) handleCreateInvestor(w http.ResponseWriter, r *http.Request) {
	var input domain.Investor

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Email == nil || input.Name == nil || input.InvestorType == nil {
		http.Error(w, "name, email, and investor_type are required", http.StatusBadRequest)
		return
	}

	if *input.Email == "" || *input.Name == "" || *input.InvestorType == "" {
		http.Error(w, "name, email, and investor_type cannot be empty", http.StatusBadRequest)
		return
	}

	newInvestor, err := h.service.Investors.CreateInvestor(input)
	if err != nil {
		http.Error(w, "could not create investor", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, newInvestor)
}

// List all investments for a specific fund
func (h *Handler) handleGetFundInvestments(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("fund_id")
	fundID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid fund UUID", http.StatusBadRequest)
		return
	}

	investments, err := h.service.Investments.GetInvestmentsByFund(fundID)
	if err != nil {
		http.Error(w, "failed to retrieve investments", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusOK, investments)
}

// Create a new investment to fund
func (h *Handler) handleCreateInvestment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("fund_id")
	fundID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid fund UUID", http.StatusBadRequest)
		return
	}

	var input domain.Investment
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.InvestorID == uuid.Nil {
		http.Error(w, "investor_id is required", http.StatusBadRequest)
		return
	}

	if input.AmountUSD == nil || input.InvestmentDate == nil {
		http.Error(w, "amount_usd and investment_date are required", http.StatusBadRequest)
		return
	}

	if *input.AmountUSD <= 0 {
		http.Error(w, "amount_usd must be greater than zero", http.StatusBadRequest)
		return
	}

	input.FundID = fundID

	newInv, err := h.service.Investments.CreateInvestment(input)
	if err != nil {
		http.Error(w, "could not record investment", http.StatusInternalServerError)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, newInv)
}
