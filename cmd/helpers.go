package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Quote struct {
	Quotation string `json:"quotation"`
	Author    string `json:"author"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	payload := ErrorResponse{}
	payload.Error = msg
	log.Println(msg)
	respondWithJSON(w, code, payload)
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(data)
}
