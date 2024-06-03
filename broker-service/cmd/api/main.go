package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	Rabbit amqp091.Connection
}

func main() {
	rabbitmqConn, err := connectRabbitMq()
	if err != nil {
		panic(err)
	}

	defer rabbitmqConn.Close()
	app := Config{
		Rabbit: *rabbitmqConn,
	}

	log.Printf("Started broker service on port %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connectRabbitMq() (*amqp091.Connection, error) {
	var counts int64
	var timeoff = 1 * time.Second
	var conn *amqp091.Connection

	for {
		c, err := amqp091.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("Could not connect to rabbitmq, retrying")
			counts += 1
		} else {
			fmt.Println("Connected to rabbit mq")
			conn = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		timeoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Waiting for %s seconds and retrying...", timeoff)
		time.Sleep(timeoff)
		continue
	}

	return conn, nil
}
