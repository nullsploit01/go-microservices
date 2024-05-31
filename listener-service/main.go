package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/nullsploit01/go-microservices/listener/event"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitmqConn, err := connect()
	if err != nil {
		panic(err)
	}

	defer rabbitmqConn.Close()

	fmt.Println("listening and consuming rabbit mq messages")

	consumer, err := event.NewConsumer(rabbitmqConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

}

func connect() (*amqp091.Connection, error) {
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
