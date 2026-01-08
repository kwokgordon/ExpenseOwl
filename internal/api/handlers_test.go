package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/kwokgordon/expenseowl/internal/config"
	"github.com/kwokgordon/expenseowl/internal/storage"
)

func TestAddExpense_ConvertsCurrency(t *testing.T) {
	// Setup temp storage and config
	tmp := t.TempDir()
	cfg := config.NewConfig(tmp)
	// set exchange rate for jpy to 0.01 CAD
	if err := cfg.UpdateExchangeRates(map[string]float64{"jpy": 0.01}); err != nil {
		t.Fatalf("failed to update exchange rates: %v", err)
	}
	st, err := storage.New(filepath.Join(tmp, "expenses.json"))
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	h := NewHandler(st, cfg)

	reqBody := map[string]interface{}{
		"name":           "Sushi",
		"category":       "Travel",
		"currency":       "jpy",
		"currencyAmount": 100.0,
		"date":           "2026-01-07T12:00:00Z",
	}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/expense", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.AddExpense(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	exps, err := st.GetAllExpenses()
	if err != nil {
		t.Fatalf("failed to read expenses: %v", err)
	}
	if len(exps) != 1 {
		t.Fatalf("expected 1 expense, got %d", len(exps))
	}
	e := exps[0]
	if e.Currency != "jpy" {
		t.Fatalf("expected currency jpy, got %s", e.Currency)
	}
	if e.CurrencyAmount != 100.0 {
		t.Fatalf("expected currencyAmount 100, got %v", e.CurrencyAmount)
	}
	// 100 * 0.01 = 1.0 CAD
	if e.Amount != 1.0 {
		t.Fatalf("expected CAD amount 1.0, got %v", e.Amount)
	}
}

func TestAddExpense_ManualCADOverride(t *testing.T) {
	// Setup temp storage and config
	tmp := t.TempDir()
	cfg := config.NewConfig(tmp)
	if err := cfg.UpdateExchangeRates(map[string]float64{"jpy": 0.01}); err != nil {
		t.Fatalf("failed to update exchange rates: %v", err)
	}
	st, err := storage.New(filepath.Join(tmp, "expenses.json"))
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	h := NewHandler(st, cfg)

	// Provide an explicit amount (CAD) that should override currency conversion
	reqBody := map[string]interface{}{
		"name":           "Dinner",
		"category":       "Travel",
		"currency":       "jpy",
		"currencyAmount": 100.0,
		"amount":         50.0,
		"date":           "2026-01-07T12:00:00Z",
	}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/expense", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.AddExpense(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	exps, err := st.GetAllExpenses()
	if err != nil {
		t.Fatalf("failed to read expenses: %v", err)
	}
	if len(exps) != 1 {
		t.Fatalf("expected 1 expense, got %d", len(exps))
	}
	e := exps[0]
	// Manual CAD amount 50 should be used despite currency fields
	if e.Amount != 50.0 {
		t.Fatalf("expected CAD amount 50.0 override, got %v", e.Amount)
	}
}

func TestReassignCategory(t *testing.T) {
	tmp := t.TempDir()
	cfg := config.NewConfig(tmp)
	// initial categories
	cfg.UpdateCategories([]string{"Old","New"})
	st, err := storage.New(filepath.Join(tmp, "expenses.json"))
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}
	// add two expenses in Old
	e1 := &config.Expense{ID: "1", Name: "One", Category: "Old", Amount: 10.0, Date: time.Now().UTC()}
	e2 := &config.Expense{ID: "2", Name: "Two", Category: "Old", Amount: 20.0, Date: time.Now().UTC()}
	if err := st.SaveExpense(e1); err != nil {
		t.Fatalf("failed to save e1: %v", err)
	}
	if err := st.SaveExpense(e2); err != nil {
		t.Fatalf("failed to save e2: %v", err)
	}
	h := NewHandler(st, cfg)
	payload := map[string]string{"old":"Old","new":"New"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/categories/delete", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ReassignCategory(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	exps, err := st.GetAllExpenses()
	if err != nil {
		t.Fatalf("failed to read expenses: %v", err)
	}
	for _, e := range exps {
		if e.Category != "New" {
			t.Fatalf("expected category New, got %s", e.Category)
		}
	}
	// config categories should not include Old
	for _, c := range cfg.Categories {
		if c == "Old" {
			t.Fatalf("expected Old to be removed from categories")
		}
	}
}
