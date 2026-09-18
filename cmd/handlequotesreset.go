package main

import "net/http"

func (apiCfg *APIConfig) handleQuotesReset(w http.ResponseWriter, r *http.Request) {
	err := apiCfg.dbQueries.QuotesReset(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)

}
