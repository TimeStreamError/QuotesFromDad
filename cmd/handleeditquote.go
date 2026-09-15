package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *APIConfig) handleEditQuote(w http.ResponseWriter, r *http.Request) {
	quoteIDString := chi.URLParam(r, "quote-id")
	quoteID, err := uuid.Parse(quoteIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	quote, err := apiCfg.dbQueries.EditQuote(r.Context(), quoteID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, quote)
}
