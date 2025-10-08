package main

import (
	//"fmt"
	"net/http"
)

func launchServer(cfg *apiConfig) {
	mux := http.NewServeMux()
	mux.Handle("/app/", 
	cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./pages")))))
	mux.Handle("/app/pages/", 
	cfg.middlewareMetricsInc(http.StripPrefix("/app/pages/", http.FileServer(http.Dir("./pages")))))
	mux.Handle("/app/assets/",
	cfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets")))))
	mux.HandleFunc("GET /admin/metrics", cfg.middlewarePrintMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.middlewareResetMetrics)
	mux.HandleFunc("GET /api/healthz", muxHandler)

	s := &http.Server{
		Handler:	mux,
		Addr:		":8080",
	}

	s.ListenAndServe()
}

func muxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}