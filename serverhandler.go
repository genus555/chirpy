package main

import (
	//"fmt"
	"net/http"
	"encoding/json"
	"log"
)

type PostRequest struct {
		Body 		string 	`json:"body"`
		CleanedBody string	`json:"cleaned_body"`
		Error 		string 	`json:"error"`
		Valid 		bool 	`json:"valid"`
		Email		string	`json:"email"`
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

func validateChirp(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	req := PostRequest{}
	err := d.Decode(&req)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	checkBadWords(&req)

	if len(req.CleanedBody) <= 140 {
		req.Valid = true
		w.WriteHeader(200)
	} else if len(req.CleanedBody) > 140 {
		req.Valid = false
		req.Error = "Chirp is too long"
		w.WriteHeader(400)
	} else {
		req.Valid = false
		req.Error = "Something went wrong"
		w.WriteHeader(500)
	}

	data, err := EncodeJSON(&req)
	if err != nil {
		log.Printf("Error encoding: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}