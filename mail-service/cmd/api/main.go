package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct{}

const webPort = "80"

func main() {
	app := Config{}
	log.Println("Started mail service on port", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
