package main

import(
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"strings"
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
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`<html>
  	<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  	</body>
	</html>`, hits)))
}

func (cfg *apiConfig) middlewareResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(int32(0))
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
	}
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func EncodeJSON(r interface{}) ([]byte, error) {
	d, err := json.Marshal(r)
	if err != nil {return d, err}

	return d, nil
}

func checkBadWords(r *PostRequest) {
	nonoWords := [3]string{"kerfuffle", "sharbert", "fornax"}

	b := strings.Split(r.Body, " ")
	for i, word := range b {
		for _, badWord := range nonoWords {
			if strings.ToLower(word) == badWord {
				b[i] = "****"
			}
		}
	}
	r.CleanedBody = strings.Join(b, " ")
}

func recievePostRequest(w http.ResponseWriter, r *http.Request) (PostRequest, error) {
	d := json.NewDecoder(r.Body)
	req := PostRequest{}
	err := d.Decode(&req)
	if err != nil {return PostRequest{}, err}

	return req, nil
}