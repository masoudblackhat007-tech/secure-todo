package main

import (
	"log"
	"net/http"

	"github.com/masoudblackhat007-tech/secure-todo/internal/config"
	apphttp "github.com/masoudblackhat007-tech/secure-todo/internal/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mux := apphttp.NewMux()

	log.Printf("starting secure-todo api on :%s (env=%s)\n", cfg.Server.Port, cfg.Server.Env)

	addr := ":" + cfg.Server.Port

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
