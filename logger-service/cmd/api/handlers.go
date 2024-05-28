package main

import (
	"net/http"

	"github.com/nullsploit01/go-microservices/logger/data"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) CreateLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload Payload

	err := app.readJson(w, r, requestPayload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	e := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(e)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	res := Response{
		Error:   false,
		Message: "logged",
	}

	app.writeJson(w, http.StatusAccepted, res)
}
