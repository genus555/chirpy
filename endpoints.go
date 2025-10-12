package main

import (
	"net/http"
)

func SetEndPoints(mux *http.ServeMux, cfg *apiConfig) {
	mux.Handle("/app/", 
	cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./pages")))))
	mux.Handle("/app/pages/", 
	cfg.middlewareMetricsInc(http.StripPrefix("/app/pages/", http.FileServer(http.Dir("./pages")))))
	mux.Handle("/app/assets/",
	cfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets")))))
	mux.HandleFunc("GET /admin/metrics", cfg.middlewarePrintMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.middlewareResetMetrics)
	mux.HandleFunc("GET /api/healthz", muxHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("POST /api/users", cfg.CreateUser)
}