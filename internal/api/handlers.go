package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kwokgordon/expenseowl/internal/config"
	"github.com/kwokgordon/expenseowl/internal/storage"
	"github.com/kwokgordon/expenseowl/internal/web"
)

type Handler struct {
	storage storage.Storage
	config  *config.Config
}

func NewHandler(s storage.Storage, cfg *config.Config) *Handler {
	return &Handler{
		storage: s,
		config:  cfg,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ExpenseRequest struct {
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Amount         float64   `json:"amount"` // Optional: CAD amount
	Currency       string    `json:"currency,omitempty"`
	CurrencyAmount float64   `json:"currencyAmount,omitempty"`
	Date           time.Time `json:"date"`
	Description    string    `json:"description,omitempty"`
}

type ConfigResponse struct {
	Categories    []string           `json:"categories"`
	Currency      string             `json:"currency"`
	StartDate     int                `json:"startDate"`
	Budgets       map[string]float64 `json:"budgets"`
	ExchangeRates map[string]float64 `json:"exchangeRates"`
}

func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	response := ConfigResponse{
		Categories:    h.config.Categories,
		Currency:      h.config.Currency,
		StartDate:     h.config.StartDate,
		Budgets:       h.config.Budgets,
		ExchangeRates: h.config.ExchangeRates,
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) EditCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	var categories []string
	if err := json.NewDecoder(r.Body).Decode(&categories); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	h.config.UpdateCategories(categories)
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
	log.Println("HTTP: Updated categories")
}

func (h *Handler) EditCurrency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	var currency string
	if err := json.NewDecoder(r.Body).Decode(&currency); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	h.config.UpdateCurrency(currency)
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
	log.Println("HTTP: Updated currency")
}

func (h *Handler) EditStartDate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	var startDate int
	if err := json.NewDecoder(r.Body).Decode(&startDate); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	h.config.UpdateStartDate(startDate)
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
	log.Println("HTTP: Updated start date")
}

func (h *Handler) EditBudgets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	var budgets map[string]float64
	if err := json.NewDecoder(r.Body).Decode(&budgets); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	h.config.UpdateBudgets(budgets)
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
	log.Println("HTTP: Updated budgets")
}

func (h *Handler) EditExchangeRates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	var rates map[string]float64
	if err := json.NewDecoder(r.Body).Decode(&rates); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	h.config.UpdateExchangeRates(rates)
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
	log.Println("HTTP: Updated exchange rates")
}

func (h *Handler) ReassignCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	var payload struct {
		Old string `json:"old"`
		New string `json:"new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	if payload.Old == "" || payload.New == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Old and New category required"})
		return
	}
	// Reassign in storage
	if err := h.storage.ReassignCategory(payload.Old, payload.New); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to reassign category"})
		log.Printf("HTTP ERROR: Failed to reassign category: %v\n", err)
		return
	}
	// Update categories list by removing old
	newCats := make([]string, 0, len(h.config.Categories))
	foundNew := false
	for _, c := range h.config.Categories {
		if c == payload.Old {
			continue
		}
		if c == payload.New {
			foundNew = true
		}
		newCats = append(newCats, c)
	}
	if !foundNew {
		newCats = append(newCats, payload.New)
	}
	h.config.UpdateCategories(newCats)
	// Migrate budgets (merge old into new)
	if h.config.Budgets != nil {
		oldB, ok := h.config.Budgets[payload.Old]
		if ok && oldB > 0 {
			h.config.Budgets[payload.New] = h.config.Budgets[payload.New] + oldB
			delete(h.config.Budgets, payload.Old)
			h.config.UpdateBudgets(h.config.Budgets)
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
	log.Printf("HTTP: Reassigned category %s -> %s\n", payload.Old, payload.New)
}

func (h *Handler) AddExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	var req ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	if !req.Date.IsZero() {
		req.Date = req.Date.UTC()
	}
	// Compute canonical Amount in CAD: prefer explicit CAD amount (manual override);
	// otherwise compute from currency fields if present
	var computedAmount float64
	if req.Amount > 0 {
		computedAmount = req.Amount // user override
	} else if req.Currency != "" && req.CurrencyAmount > 0 {
		rate, exists := h.config.ExchangeRates[strings.ToLower(req.Currency)]
		if !exists {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Missing exchange rate for currency - please set it in settings"})
			log.Printf("HTTP ERROR: Missing exchange rate for %s\n", req.Currency)
			return
		}
		computedAmount = req.CurrencyAmount * rate
	} else {
		computedAmount = req.Amount // could be 0
	}
	       expense := &config.Expense{
		       Name:           req.Name,
		       Category:       req.Category,
		       Amount:         computedAmount,
		       Currency:       req.Currency,
		       CurrencyAmount: req.CurrencyAmount,
		       Date:           req.Date,
		       Description:    req.Description,
	       }
	if err := expense.Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		log.Printf("HTTP ERROR: Failed to validate expense: %v\n", err)
		return
	}
	if err := h.storage.SaveExpense(expense); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to save expense"})
		log.Printf("HTTP ERROR: Failed to save expense: %v\n", err)
		return
	}
	writeJSON(w, http.StatusOK, expense)
}

func (h *Handler) EditExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "ID parameter is required"})
		log.Println("HTTP ERROR: ID parameter is required")
		return
	}
	var req ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		log.Printf("HTTP ERROR: Failed to decode request body: %v\n", err)
		return
	}
	if !req.Date.IsZero() {
		req.Date = req.Date.UTC()
	}
	// Compute canonical Amount in CAD for edit: prefer explicit CAD amount (manual override);
	// otherwise compute from currency fields if present
	var computedAmount float64
	if req.Amount > 0 {
		computedAmount = req.Amount // user override
	} else if req.Currency != "" && req.CurrencyAmount > 0 {
		rate, exists := h.config.ExchangeRates[strings.ToLower(req.Currency)]
		if !exists {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Missing exchange rate for currency - please set it in settings"})
			log.Printf("HTTP ERROR: Missing exchange rate for %s\n", req.Currency)
			return
		}
		computedAmount = req.CurrencyAmount * rate
	} else {
		computedAmount = req.Amount // could be 0
	}
	       expense := &config.Expense{
		       ID:             id,
		       Name:           req.Name,
		       Category:       req.Category,
		       Amount:         computedAmount,
		       Currency:       req.Currency,
		       CurrencyAmount: req.CurrencyAmount,
		       Date:           req.Date,
		       Description:    req.Description,
	       }
	if err := expense.Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		log.Printf("HTTP ERROR: Failed to validate expense: %v\n", err)
		return
	}
	if err := h.storage.EditExpense(expense); err != nil {
		if err == storage.ErrExpenseNotFound {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "Expense not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to edit expense"})
		log.Printf("HTTP ERROR: Failed to edit expense: %v\n", err)
		return
	}
	writeJSON(w, http.StatusOK, expense)
	log.Printf("HTTP: Edited expense with ID %s\n", id)
}

func (h *Handler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	expenses, err := h.storage.GetAllExpenses()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve expenses"})
		log.Printf("HTTP ERROR: Failed to retrieve expenses: %v\n", err)
		return
	}
	writeJSON(w, http.StatusOK, expenses)
}

func (h *Handler) ServeTableView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := web.ServeTemplate(w, "table.html"); err != nil {
		http.Error(w, "Failed to serve template", http.StatusInternalServerError)
		log.Printf("HTTP ERROR: Failed to serve template: %v\n", err)
		return
	}
}

func (h *Handler) ServeSettingsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := web.ServeTemplate(w, "settings.html"); err != nil {
		http.Error(w, "Failed to serve template", http.StatusInternalServerError)
		log.Printf("HTTP ERROR: Failed to serve template: %v\n", err)
		return
	}
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "ID parameter is required"})
		log.Println("HTTP ERROR: ID parameter is required")
		return
	}
	if err := h.storage.DeleteExpense(id); err != nil {
		if err == storage.ErrExpenseNotFound {
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "Expense not found"})
			log.Printf("HTTP ERROR: Expense not found: %v\n", err)
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete expense"})
		log.Printf("HTTP ERROR: Failed to delete expense: %v\n", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
	log.Printf("HTTP: Deleted expense with ID %s\n", id)
}

// Static Handler
func (h *Handler) ServeStaticFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("HTTP ERROR: Method not allowed")
		return
	}
	if err := web.ServeStatic(w, r.URL.Path); err != nil {
		http.Error(w, "Failed to serve static file", http.StatusInternalServerError)
		log.Printf("HTTP ERROR: Failed to serve static file %s: %v\n", r.URL.Path, err)
		return
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
