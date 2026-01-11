package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	ServerPort    string
	StoragePath   string
	Categories    []string
	Currency      string
	StartDate     int
	Budgets       map[string]float64
	ExchangeRates map[string]float64
	mu            sync.RWMutex
}

type FileConfig struct {
	Categories    []string           `json:"categories"`
	Currency      string             `json:"currency"`
	StartDate     int                `json:"startDate"`
	Budgets       map[string]float64 `json:"budgets,omitempty"`
	ExchangeRates map[string]float64 `json:"exchangeRates,omitempty"`
}

var defaultCategories = []string{
	"Air Fare",
	"Hotel",
	"Food",
	"Transit",
	"Gordon",
	"Angela",
	"Evie",
	"eSim",
}

var currencySymbols = map[string]string{
	"cad": "C$",   // Canadian Dollar
	"usd": "$",    // US Dollar
}

type Expense struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Amount         float64   `json:"amount"`          // Amount in CAD (stored canonical)
	Currency       string    `json:"currency,omitempty"`       // Currency code used when purchasing (e.g., JPY)
	CurrencyAmount float64   `json:"currencyAmount,omitempty"` // Original amount in the transactional currency
	Date           time.Time `json:"date"`
}

func (e *Expense) Validate() error {
	if e.Name == "" {
		return errors.New("expense name is required")
	}
	if e.Category == "" {
		return errors.New("category is required")
	}
	if e.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

func NewConfig(dataPath string) *Config {
	finalPath := ""
	if dataPath == "data" {
		finalPath = filepath.Join(".", "data")
	} else {
		finalPath = filepath.Clean(dataPath)
	}
	if err := os.MkdirAll(finalPath, 0755); err != nil {
		log.Printf("Error creating data directory: %v", err)
	}
	log.Printf("Using data directory: %s\n", finalPath)
	cfg := &Config{
		ServerPort:    "8080",
		StoragePath:   finalPath,
		Categories:    defaultCategories,
		StartDate:     1,
		Currency:      "C$", // Default to CAD symbol
		Budgets:       make(map[string]float64),
		ExchangeRates: map[string]float64{"cad": 1.0},
	}
	configPath := filepath.Join(finalPath, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("Configuration file not found, creating default configuration")
		if envCategories := os.Getenv("EXPENSE_CATEGORIES"); envCategories != "" {
			categories := strings.Split(envCategories, ",")
			for i := range categories {
				categories[i] = strings.TrimSpace(categories[i])
			}
			cfg.Categories = categories
			log.Println("Using custom categories from environment variables")
		}
		if envCurrency := strings.ToLower(os.Getenv("CURRENCY")); envCurrency != "" {
			if symbol, exists := currencySymbols[envCurrency]; exists {
				cfg.Currency = symbol
			}
			log.Println("Using custom currency from environment variables")
		}
		if envStartDate := strings.ToLower(os.Getenv("START_DATE")); envStartDate != "" {
			startDate, err := strconv.Atoi(envStartDate)
			if err != nil {
				log.Println("START_DATE is not a number, using default (1)")
			} else {
				cfg.StartDate = startDate
				log.Println("using custom start date from environment variables")
			}
		}
	} else if fileConfig, err := loadConfigFile(configPath); err == nil {
		cfg.Categories = fileConfig.Categories
		if fileConfig.Currency != "" {
			cfg.Currency = fileConfig.Currency
		}
		cfg.StartDate = fileConfig.StartDate
		if fileConfig.Budgets != nil {
			cfg.Budgets = fileConfig.Budgets
		}
		if fileConfig.ExchangeRates != nil {
			cfg.ExchangeRates = fileConfig.ExchangeRates
		}
		log.Println("Loaded configuration from file")
	}
	cfg.SaveConfig()
	return cfg
}

func loadConfigFile(filePath string) (*FileConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var config FileConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) SaveConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	filePath := filepath.Join(c.StoragePath, "config.json")
	fileConfig := FileConfig{
		Categories:    c.Categories,
		Currency:      c.Currency,
		StartDate:     c.StartDate,
		Budgets:       c.Budgets,
		ExchangeRates: c.ExchangeRates,
	}
	data, err := json.MarshalIndent(fileConfig, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func (c *Config) UpdateCategories(categories []string) error {
	c.mu.Lock()
	c.Categories = categories
	c.mu.Unlock()
	return c.SaveConfig()
}

func (c *Config) UpdateCurrency(currencyCode string) error {
	c.mu.Lock()
	if symbol, exists := currencySymbols[strings.ToLower(currencyCode)]; exists {
		c.Currency = symbol
	} else {
		c.mu.Unlock()
		return errors.New("invalid currency code")
	}
	c.mu.Unlock()
	return c.SaveConfig()
}

func (c *Config) UpdateStartDate(startDate int) error {
	c.mu.Lock()
	c.StartDate = max(min(startDate, 31), 1)
	c.mu.Unlock()
	return c.SaveConfig()
}

func (c *Config) UpdateBudgets(budgets map[string]float64) error {
	c.mu.Lock()
	c.Budgets = budgets
	c.mu.Unlock()
	return c.SaveConfig()
}

func (c *Config) UpdateExchangeRates(rates map[string]float64) error {
	c.mu.Lock()
	c.ExchangeRates = rates
	c.mu.Unlock()
	return c.SaveConfig()
}
