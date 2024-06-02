package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"

	"github.com/nullsploit01/go-microservices/broker/event"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"Data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := Response{
		Error:   false,
		Message: "Broker says what?",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	if err := app.readJson(w, r, &requestPayload); err != nil {
		app.errorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	case "log":
		app.logItemViaRPC(w, requestPayload.Log)

	case "mail":
		app.sendMail(w, requestPayload.Mail)

	default:
		app.errorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	jsonData, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	var jsonResponseFromService Response
	err = json.NewDecoder(response.Body).Decode(&jsonResponseFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if response.StatusCode != http.StatusCreated {
		app.errorJson(w, errors.New(jsonResponseFromService.Message))
		return
	}

	var payload Response
	payload.Error = false
	payload.Message = jsonResponseFromService.Message
	payload.Data = jsonResponseFromService.Data

	app.writeJson(w, http.StatusCreated, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()

	var jsonResponseFromService Response
	err = json.NewDecoder(response.Body).Decode(&jsonResponseFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New(jsonResponseFromService.Message), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New(jsonResponseFromService.Message))
		return
	}

	if jsonResponseFromService.Error {
		app.errorJson(w, errors.New(jsonResponseFromService.Message), http.StatusUnauthorized)
		return
	}

	var payload Response
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = jsonResponseFromService.Data

	app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	jsonData, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	defer response.Body.Close()

	var jsonResponseFromService Response
	err = json.NewDecoder(response.Body).Decode(&jsonResponseFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if response.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New(jsonResponseFromService.Message))
		return
	}

	var payload Response
	payload.Error = false
	payload.Message = jsonResponseFromService.Message
	payload.Data = jsonResponseFromService.Data

	app.writeJson(w, http.StatusCreated, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var response Response
	response.Error = false
	response.Message = "event logged!"

	app.writeJson(w, http.StatusOK, response)
}

func (app *Config) pushToQueue(name, message string) error {
	emitter, err := event.NewEventEmitter(&app.Rabbit)

	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: message,
	}

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), "log.INFO")

	return err
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", payload, &result)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	resp := Response{
		Error:   false,
		Message: result,
	}

	app.writeJson(w, http.StatusOK, resp)
}
