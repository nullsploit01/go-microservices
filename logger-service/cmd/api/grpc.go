package main

import (
	"context"

	"github.com/nullsploit01/go-microservices/logger/data"
	"github.com/nullsploit01/go-microservices/logger/logs"
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
