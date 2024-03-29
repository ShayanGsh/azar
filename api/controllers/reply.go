package controllers

import (
	"encoding/json"
	"net/http"
)

type ReplyMessage struct {
	Success bool
	Message string
	Status  int
}

func Reply(rw http.ResponseWriter, message ReplyMessage) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(message.Status)
	json.NewEncoder(rw).Encode(message)
}

func ReplyError(rw http.ResponseWriter, err error, status int) {
	Reply(rw, ReplyMessage{
		Success: false,
		Message: err.Error(),
		Status:  status,
	})
}

func ReplySuccess(rw http.ResponseWriter, message string, status ...int) {
	if len(status) == 0 {
		status = append(status, http.StatusOK)
	}
	Reply(rw, ReplyMessage{
		Success: true,
		Message: message,
		Status: status[0],
	})
}