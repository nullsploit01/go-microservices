package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitmqConn, err := connect()
	if err != nil {
		panic(err)
	}

	defer rabbitmqConn.Close()
}

func connect() (*amqp091.Connection, error) {
	var counts int64
	var timeoff = 1 * time.Second
	var conn *amqp091.Connection

	for {
		c, err := amqp091.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("Could not connect to rabbitmq, retrying")
			counts += 1
		} else {
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
