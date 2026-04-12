package main

import (
	"log"
	"net/http"
	"os"

	"hybrid-app/backend/internal/app"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("go backend listening on %s", addr)
	if err := http.ListenAndServe(addr, application.Router()); err != nil {
		log.Fatal(err)
	}
}
