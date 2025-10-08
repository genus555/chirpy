package main

import(
	"net/http"
	"log"
	"fmt"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		_ = cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	}

	return http.HandlerFunc(f)
}

func (cfg *apiConfig) middlewarePrintMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	fmt.Println(string(hits))
	w.Write([]byte(fmt.Sprintf("Hits: %d\n", hits)))
}

func (cfg *apiConfig) middlewareResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(int32(0))
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}