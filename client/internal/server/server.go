package server

import (
	"log"
	"lzhuk/clients/internal/router"
	"lzhuk/clients/pkg/config"
	"net/http"
)

func StartServer(cfg config.Config) {
	server := router.Router(cfg)
	err := server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
