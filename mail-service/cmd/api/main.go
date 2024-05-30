package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

func main() {
	m, err := initMailer()
	if err != nil {
		log.Panic(err)
	}

	app := Config{
		Mailer: m,
	}

	log.Println("Started mail service on port", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func initMailer() (Mail, error) {
	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return Mail{}, err
	}

	mailer := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_NAME"),
		FromAddress: os.Getenv("MAIL_ADDRESS"),
	}

	return mailer, nil
}
