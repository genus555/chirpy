package main

import (
	//"fmt"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}
func main() {
	var serverHits atomic.Int32
	initializeHits(&serverHits)
	cfg := apiConfig{
		fileserverHits: serverHits,
	}
	launchServer(&cfg);
}

func initializeHits(hits *atomic.Int32) {
	hits.Store(int32(0))
}