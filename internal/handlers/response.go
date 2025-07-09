package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Token   string      `json:"token,omitempty"`
}

func responseSuccess(w http.ResponseWriter, T any, msg string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(Response{
		Success: true,
		Data:    T,
		Message: msg,
	})

	if err != nil {
		return fmt.Errorf("an error occurred while responding success message : %v", err)
	}

	return nil
}
