package main

import (
	"fmt"
	"net/http"
	
)

func main() {
	mux := http.NewServeMux()

	s := http.Server{
		Handler:	mux,
		Addr:		":8080",
	}

	s.ListenAndServe()

	fmt.Println(mux, s.Addr)
}