package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/nullsploit01/go-microservices/authentication/data"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}
