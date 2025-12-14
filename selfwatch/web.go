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

func (ws *WebServer) handleHourly(w http.ResponseWriter, r *http.Request) {
	counts, err := ws.Storage.HourlyCounts(24, 0) // No day offset for hourly data
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
