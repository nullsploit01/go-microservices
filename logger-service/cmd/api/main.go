package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nullsploit01/go-microservices/logger/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

type Config struct {
	Models data.Models
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	client, err := connectMongoDB(ctx)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	go app.server()
	log.Printf("Logger service running on port %s\n", webPort)
}

func (app *Config) server() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func connectMongoDB(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Println("Error connecting MongoDB", err)
		return nil, err
	}

	return c, nil
}
