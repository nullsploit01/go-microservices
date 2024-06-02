package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/nullsploit01/go-microservices/logger/data"
	"github.com/nullsploit01/go-microservices/logger/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	ip := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: ip.Name,
		Data: ip.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "Failed",
		}

		return res, err
	}

	res := &logs.LogResponse{
		Result: "logged successfully",
	}

	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc %v", err)
	}

	server := grpc.NewServer()
	logs.RegisterLogServiceServer(server, &LogServer{Models: app.Models})

	log.Println("gRPC Server started on port", gRpcPort)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server %v", err)
	}
}
