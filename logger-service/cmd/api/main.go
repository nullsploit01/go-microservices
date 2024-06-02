package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
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

var client *mongo.Client

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	c, err := connectMongoDB(ctx)
	if err != nil {
		log.Panic(err)
	}

	client = c

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Logger service running on port %s\n", webPort)

	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC Server on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", rpcPort))

	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConnection, err := listen.Accept()
		if err != nil {
			fmt.Println("error accepting tcp connection", err)
			continue
		}

		go rpc.ServeConn(rpcConnection)
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

	log.Println("Connected to mongo!")

	return c, nil
}
