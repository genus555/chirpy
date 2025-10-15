package main

import (
	//"fmt"
	"net/http"
	"github.com/google/uuid"
)

type PostRequest struct {
		Body 		string 		`json:"body"`
		CleanedBody string		`json:"cleaned_body"`
		Error 		string 		`json:"error"`
		Valid 		bool 		`json:"valid"`
		Email		string		`json:"email"`
		UserID		uuid.UUID	`json:"user_id"`
		Password	string		`json:"password"`
	}

func launchServer(cfg *apiConfig) {
	mux := http.NewServeMux()
	SetEndPoints(mux, cfg)

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