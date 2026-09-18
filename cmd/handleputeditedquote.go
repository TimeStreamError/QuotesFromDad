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

	// get what we want to put in the db from the request
	decoder := json.NewDecoder(r.Body)

	type ReturnedQuote struct {
		Quote  string   `json:"quotation"`
		Author string   `json:"author"`
		ID     string   `json:"id"`
		Tags   []string `json:"tags"`
	}

	data := ReturnedQuote{}
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Now we have the data, need to open a transaction and do three queries
	// First set the new quote and author, then delete all previous tags, then set the new tags

	// set up the transaction
	tx, err := apiCfg.DB.Begin()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// After we are done, rollback. If we committed, we roll back nothing, otherwise we restore the db
	defer tx.Rollback()

	// now we set up to use the transaction with the queries we need.
	qtx := apiCfg.dbQueries.WithTx(tx)

	// First up, put the edited quote into the quotes table
	params := database.PutEditedQuoteParams{}
	params.Author.Valid = false
	if data.Author != "" {
		params.Author.String = data.Author
		params.Author.Valid = true
	}
	params.UpdatedAt = time.Now()
	params.Quote = data.Quote
	params.ID, err = uuid.Parse(data.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseQuote, err := qtx.PutEditedQuote(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Next, delete all the tag links associated with this quote
	err = qtx.DeleteQuoteTags(r.Context(), params.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseJson := struct {
		Quote string `json:"quote"`
	}{Quote: responseQuote.Quote}
	respondWithJSON(w, http.StatusOK, responseJson)
}
