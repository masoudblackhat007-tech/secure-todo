package main

import (
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		// در همین مرحله لاگ ساده کافی است
		log.Printf("write response error: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)

	addr := ":8081"
	log.Printf("starting server on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
