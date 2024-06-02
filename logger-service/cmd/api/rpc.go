package main

import (
	"context"
	"log"
	"time"

	"github.com/nullsploit01/go-microservices/logger/data"
)

type RPCServer struct{}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(paylaod RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.Background(), data.LogEntry{
		Name:      paylaod.Name,
		Data:      paylaod.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error adding log", err)
		return err
	}

	*resp = "Processed payload via RPC " + paylaod.Name

	return nil
}
