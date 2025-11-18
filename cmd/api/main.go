package main

import (
	"log"
	"net/http"

	"github.com/masoudblackhat007-tech/secure-todo/internal/config"
	apphttp "github.com/masoudblackhat007-tech/secure-todo/internal/http"
)

func main() {
	cfg := config.Load()

	mux := apphttp.NewMux()

	log.Printf("starting secure-todo api on :%s\n", cfg.HTTPPort)

	if err := http.ListenAndServe(":"+cfg.HTTPPort, mux); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
