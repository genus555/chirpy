package main

import (
	//"fmt"
	"log"
	"net/http"
	"time"
	"database/sql"

	"github.com/google/uuid"
	"github.com/genus555/chirpy/internal/database"
	"github.com/genus555/chirpy/internal/auth"
)

type User struct {
	ID					uuid.UUID	`json:"id"`
	CreatedAt			time.Time	`json:"created_at"`
	UpdatedAt			time.Time	`json:"updated_at"`
	Email				string		`json:"email"`
	Token				string		`json:"token"`
	RefreshToken		string		`json:"refresh_token"`
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
	if req.Email == "" || req.Password == "" {
		log.Printf("Error no password or email")
		w.WriteHeader(400)
		return
	}

	h_pswrd, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	u, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:				req.Email,
		HashedPassword:		h_pswrd,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}
	user := User{
		ID:				u.ID,
		CreatedAt:		u.CreatedAt,
		UpdatedAt:		u.UpdatedAt,
		Email:			u.Email,
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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error retrieving token: %s", err)
		w.WriteHeader(401)
		return
	}

	uID, err := auth.ValidateJWT(token, cfg.ts)
	if err != nil {
		log.Printf("Token invalid: %s", err)
		w.WriteHeader(401)
		return
	}

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

	chirp, err := cfg.recieveChirp(req, r, uID)
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

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	req, err := recievePostRequest(w, r)
	//check to see if request has all necessary fields
	if err != nil {
		log.Printf("Error recieving request: %s", err)
		w.WriteHeader(500)
		return
	}
	if req.Password == "" || req.Email == "" {
		log.Printf("Missing password or email")
		w.WriteHeader(400)
		return
	}

	//set access token duration
	duration, _ := time.ParseDuration("1h")

	//get user from email
	u, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error getting user information: %s", err)
		w.WriteHeader(500)
		return
	}

	//check password
	ok, err := auth.CheckPasswordHash(req.Password, u.HashedPassword)
	if err != nil {
		log.Printf("Error checking password: %s", err)
		w.WriteHeader(500)
		return
	}

	//create token
	token, err := auth.MakeJWT(u.ID, cfg.ts, duration)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		w.WriteHeader(500)
		return
	}
	nullToken := sql.NullString{String: token, Valid: true}

	//update user with access token
	err = cfg.db.AddUserToken(r.Context(), database.AddUserTokenParams{
		ID:		u.ID,
		Token:	nullToken,
	})
	if err != nil {
		log.Printf("Error giving user token: %s", err)
		w.WriteHeader(500)
		return
	}

	//set refresh token duration
	duration, _ = time.ParseDuration("60d")

	//add refresh token to database
	r_token, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		w.WriteHeader(500)
		return
	}
	rt, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:		r_token,
		UserID:		u.ID,
		ExpiresAt:	time.Now().Add(duration),
	})

	//populate User struct
	user := User{
		ID:				u.ID,
		CreatedAt:		u.CreatedAt,
		UpdatedAt:		u.UpdatedAt,
		Email:			u.Email,
		Token:			token,
		RefreshToken:	rt.Token,
	}

	//if all is right login
	if !ok {
		log.Printf("Incorrect email or password")
		w.WriteHeader(401)
	} else {
		data, err := EncodeJSON(user)
		if err != nil {
			log.Printf("Error encoding user: %s", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	r_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %s", err)
		w.WriteHeader(500)
		return
	}
	if r_token == "" {
		log.Printf("No refresh token located: %s", err)
		w.WriteHeader(500)
		return
	}

	rt, err := cfg.db.GetRefreshToken(r.Context(), r_token)
	if err != nil {
		log.Printf("Error getting refresh token from database: %s", err)
		w.WriteHeader(401)
		return
	}

	if rt.RevokedAt.Valid {
		log.Printf("Refresh token has been revoked")
		w.WriteHeader(401)
		return
	}

	expired := time.Now().Before(rt.ExpiresAt)

	if expired {
		log.Printf("Token has expired")
		w.WriteHeader(401)
		return
	} else {
		id, err := cfg.db.GetUserIDFromRefreshToken(r.Context(), rt.Token)
		if err != nil {
			log.Printf("Error finding user from token: %s", err)
			w.WriteHeader(500)
			return
		}

		duration, _ := time.ParseDuration("1h")
		new_token, err := auth.MakeJWT(id, cfg.ts, duration)
		nullToken := sql.NullString{String: new_token, Valid: true}

		err = cfg.db.AddUserToken(r.Context(), database.AddUserTokenParams{
			ID:			id,
			Token:		nullToken,
		})
		if err != nil {
			log.Printf("Error giving user new token: %s", err)
			w.WriteHeader(500)
			return
		}

		type Token struct {
			Token	string	`json:"token"`
		}
		access_token := Token{
			Token:		new_token,
		}

		data, err := EncodeJSON(access_token)
		if err != nil {
			log.Printf("Error encoding new access token: %s", err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	r_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting refresh token: %s", err)
		w.WriteHeader(500)
		return
	}
	if r_token == "" {
		log.Printf("No refresh token located: %s", err)
		w.WriteHeader(404)
		return
	}

	revoked := sql.NullTime{
		Time:	time.Now(),
		Valid:	true,
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt:		revoked,
		Token:			r_token,
	})
	if err != nil {
		log.Printf("Error revoking token: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}

func (cfg *apiConfig) UpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		w.WriteHeader(401)
		return
	}

	nullToken := sql.NullString{String: token, Valid: true}
	u, err := cfg.db.GetUserFromToken(r.Context(), nullToken)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		w.WriteHeader(401)
		return
	}
	
	req, err := recievePostRequest(w, r)
	if err != nil {
		log.Printf("Error recieving request: %s", err)
		w.WriteHeader(500)
		return
	}
	if req.Password == "" || req.Email == "" {
		log.Printf("Missing password or email")
		w.WriteHeader(400)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing new password: %s", err)
		w.WriteHeader(500)
		return
	}

	updated_user, err := cfg.db.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		Email:				req.Email,
		HashedPassword:		hash,
		ID:					u.ID,
	})

	user := dbUserIntoUserStruct(updated_user)
	data, err := EncodeJSON(user)
	if err != nil {
		log.Printf("Error encoding updated user: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (cfg *apiConfig) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	//check for token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		w.WriteHeader(401)
		return
	}

	//get user from token
	nullToken := sql.NullString{String: token, Valid: true}
	u, err := cfg.db.GetUserFromToken(r.Context(), nullToken)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		w.WriteHeader(401)
		return
	}

	//get request's chirp_id
	id_from_path := r.PathValue("chirp_id")
	chirp_id, err := uuid.Parse(id_from_path)
	if err != nil {
		log.Printf("Error getting chirp id: %s", err)
		w.WriteHeader(500)
		return
	}

	//get chirp from chirp_id
	chirp, err := cfg.db.GetChirpByChirpId(r.Context(), chirp_id)
	if err != nil {
		log.Printf("Error getting chirp: %s", err)
		w.WriteHeader(404)
		return
	}

	if chirp.UserID != u.ID {
		log.Printf("Incorrect user")
		w.WriteHeader(403)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		log.Printf("Error deleteing chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(204)
}