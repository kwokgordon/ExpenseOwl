package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/kwokgordon/expenseowl/internal/api"
	"github.com/kwokgordon/expenseowl/internal/config"
	"github.com/kwokgordon/expenseowl/internal/storage"
	"github.com/kwokgordon/expenseowl/internal/web"
)

func runServer(dataPath string) {
	cfg := config.NewConfig(dataPath)
	storage, err := storage.New(filepath.Join(cfg.StoragePath, "expenses.json"))
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	handler := api.NewHandler(storage, cfg)
	http.HandleFunc("/categories", handler.GetCategories)
	http.HandleFunc("/categories/edit", handler.EditCategories)
	http.HandleFunc("/currency", handler.EditCurrency)
	http.HandleFunc("/startdate", handler.EditStartDate)
	http.HandleFunc("/expense", handler.AddExpense)
	http.HandleFunc("/expenses", handler.GetExpenses)
	http.HandleFunc("/expense/edit", handler.EditExpense)
	http.HandleFunc("/table", handler.ServeTableView)
	http.HandleFunc("/settings", handler.ServeSettingsPage)
	http.HandleFunc("/expense/delete", handler.DeleteExpense)
	http.HandleFunc("/export/json", handler.ExportJSON)
	http.HandleFunc("/import/csv", handler.ImportCSV)
	http.HandleFunc("/import/json", handler.ImportJSON)
	http.HandleFunc("/export/csv", handler.ExportCSV)
	http.HandleFunc("/manifest.json", handler.ServeStaticFile)
	http.HandleFunc("/sw.js", handler.ServeStaticFile)
	http.HandleFunc("/pwa/", handler.ServeStaticFile)
	http.HandleFunc("/style.css", handler.ServeStaticFile)
	http.HandleFunc("/favicon.ico", handler.ServeStaticFile)
	http.HandleFunc("/chart.min.js", handler.ServeStaticFile)
	http.HandleFunc("/fa.min.css", handler.ServeStaticFile)
	http.HandleFunc("/webfonts/", handler.ServeStaticFile)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if err := web.ServeTemplate(w, "index.html"); err != nil {
			log.Printf("HTTP ERROR: Failed to serve template: %v", err)
			http.Error(w, "Failed to serve template", http.StatusInternalServerError)
			return
		}
	})
	log.Printf("Starting server on port %s...\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func main() {
	dataPath := flag.String("data", "data", "Path to data directory")
	flag.Parse()
	runServer(*dataPath)
}
