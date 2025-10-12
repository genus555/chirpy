package main

import (
	//"fmt"
	"log"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email		string		`json:"email"`
}

func (cfg *apiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	req, err := recievePostRequest(w, r)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	u, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	user := User{
		ID:			u.ID,
		CreatedAt:	u.CreatedAt,
		UpdatedAt:	u.UpdatedAt,
		Email:		u.Email,
	}

	data, err := EncodeJSON(&user)
	if err != nil {
		log.Printf("Error encoding user: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}