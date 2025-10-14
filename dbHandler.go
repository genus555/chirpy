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

type Chirp struct {
	ID			uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Body		string		`json:"body"`
	UserID		uuid.UUID	`json:"user_id"`
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

func (cfg *apiConfig) chirps(w http.ResponseWriter, r *http.Request) {
	req, err := recievePostRequest(w, r)
	if err != nil {
		log.Printf("Error recieving request: %s", err)
		w.WriteHeader(500)
		return
	}
	checkBadWords(&req)

	if len(req.CleanedBody) <= 140 {
		req.Valid = true
	} else if len(req.CleanedBody) > 140 {
		req.Valid = false
		req.Error = "Chirp is too long"
		w.WriteHeader(400)
	} else {
		req.Valid = false
		req.Error = "Something went wrong"
		w.WriteHeader(500)
	}
	if req.Valid == false {
		data, err := EncodeJSON(req)
		if err != nil {
			log.Printf("Error encoding: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	chirp, err := cfg.recieveChirp(req, r)
	if err != nil {
		log.Printf("Error reicieving chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	data, err := EncodeJSON(chirp)
	if err != nil {
		log.Printf("Error encoding: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps from database: %s", err)
		w.WriteHeader(500)
		return
	}

	list_of_chirps := make([]Chirp, 0, 5)
	for _, c := range chirps {
		chirp := dbChirpIntoChirpStruct(c)
		list_of_chirps = append(list_of_chirps, chirp)
	}

	data, err := EncodeJSON(list_of_chirps)
	if err != nil {
		log.Printf("Error encoding data: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (cfg *apiConfig) getChirpByChirpID(w http.ResponseWriter, r *http.Request) {
	id_from_path := r.PathValue("chirp_id")
	chirp_id, err := uuid.Parse(id_from_path)
	if err != nil {
		log.Printf("Error getting chirp id: %s", err)
		w.WriteHeader(500)
		return
	}

	chirp, err := cfg.db.GetChirpByChirpId(r.Context(), chirp_id)
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		w.WriteHeader(404)
		return
	}

	c := dbChirpIntoChirpStruct(chirp)
	data, err := EncodeJSON(c)
	if err != nil {
		log.Printf("Error encoding chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}