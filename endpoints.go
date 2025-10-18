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
	
	//GET
	mux.HandleFunc("GET /admin/metrics", cfg.middlewarePrintMetrics)
	mux.HandleFunc("GET /api/healthz", muxHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", cfg.getChirpByChirpID)

	//POST
	mux.HandleFunc("POST /admin/reset", cfg.middlewareResetMetrics)
	mux.HandleFunc("POST /api/users", cfg.CreateUser)
	mux.HandleFunc("POST /api/chirps", cfg.chirps)
	mux.HandleFunc("POST /api/login", cfg.login)
	mux.HandleFunc("POST /api/refresh", cfg.refresh)
	mux.HandleFunc("POST /api/revoke", cfg.revoke)

	//PUT
	mux.HandleFunc("PUT /api/users", cfg.UpdateUser)
}