package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/timestreamerror/capstone/internal/database"
)

func (apiCfg *APIConfig) handlePutEditedQuote(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type ReturnedQuote struct {
		Quote  string `json:"quotation"`
		Author string `json:"author"`
		ID     string `json:"id"`
	}

	data := ReturnedQuote{}
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	params := database.PutEditedQuoteParams{}
	params.Author.Valid = false
	if data.Author != "" {
		params.Author.String = data.Author
		params.Author.Valid = true
	}
	params.UpdatedAt = time.Now()
	params.ID, err = uuid.Parse(data.ID)
	params.Quote = data.Quote

	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	responseQuote, err := apiCfg.dbQueries.PutEditedQuote(r.Context(), params)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	responseJson := struct {
		Quote string `json:"quote"`
	}{Quote: responseQuote.Quote}
	respondWithJSON(w, http.StatusOK, responseJson)
}
