package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := Response{
		Error:   false,
		Message: "Broker says what?",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}
