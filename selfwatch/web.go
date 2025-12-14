package selfwatch

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"time"
)

//go:embed web/*
var webAssets embed.FS

type WebServer struct {
	Storage    *WatchStorage
	Config     *config
	ListenAddr string
}

func NewWebServer(storage *WatchStorage, cfg *config, addr string) *WebServer {
	return &WebServer{
		Storage:    storage,
		Config:     cfg,
		ListenAddr: addr,
	}
}

func (ws *WebServer) Start() error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/hourly", ws.handleHourly)
	mux.HandleFunc("/api/daily", ws.handleDaily)
	mux.HandleFunc("/api/yearly", ws.handleYearly)

	// Static files
	webFS, err := fs.Sub(webAssets, "web")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(webFS)))

	log.Printf("Starting web dashboard at http://%s", ws.ListenAddr)
	return http.ListenAndServe(ws.ListenAddr, mux)
}

func isValidDateFormat(date string) bool {
	if len(date) != 10 {
		return false
	}
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func (ws *WebServer) handleHourly(w http.ResponseWriter, r *http.Request) {
	var counts []HourlyCount
	var err error

	// Check for specific date first (e.g., ?date=2024-12-10)
	if dateParam := r.URL.Query().Get("date"); dateParam != "" {
		if !isValidDateFormat(dateParam) {
			http.Error(w, "Invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		counts, err = ws.Storage.HourlyCountsForDate(dateParam)
	} else {
		// Fall back to offset-based query
		offset := 0
		if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
			if parsed, err := strconv.Atoi(offsetParam); err == nil && parsed >= 0 {
				offset = parsed
			}
		}
		counts, err = ws.Storage.HourlyCounts(24, offset)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}

func (ws *WebServer) handleDaily(w http.ResponseWriter, r *http.Request) {
	counts, err := ws.Storage.DailyCounts(30, ws.Config.NewDayHour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}

func (ws *WebServer) handleYearly(w http.ResponseWriter, r *http.Request) {
	year := time.Now().Year()
	if yearParam := r.URL.Query().Get("year"); yearParam != "" {
		if parsed, err := strconv.Atoi(yearParam); err == nil && parsed >= 1970 && parsed <= year {
			year = parsed
		}
	}
	counts, err := ws.Storage.YearlyCounts(year, ws.Config.NewDayHour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}
