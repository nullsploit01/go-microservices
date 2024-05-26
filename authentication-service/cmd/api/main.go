package main

import (
	"database/sql"

	"github.com/nullsploit01/go-microservices/authentication/data"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {

}
