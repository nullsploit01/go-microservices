package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	if err := app.readJson(w, r, &requestPayload); err != nil {
		app.errorJson(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err := app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := Response{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	err = app.logRequest("mail", fmt.Sprintf("%s send email to %s", requestPayload.From, requestPayload.To))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return nil
}
