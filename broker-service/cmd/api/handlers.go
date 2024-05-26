package main

import (
	"net/http"

	"github.com/nullsploit01/go-microservices/broker/common/util"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := Response{
		Error:   false,
		Message: "Broker says what?",
	}

	util.WriteJSON(w, http.StatusOK, payload)
}
