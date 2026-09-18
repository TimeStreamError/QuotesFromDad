package main

import (
	"html/template"
	"net/http"

	"github.com/timestreamerror/capstone/internal/auth"
)

func (apiCfg *APIConfig) handleGetAllQuotes(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	_, err = auth.ValidateJWT(token, apiCfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	records, err := apiCfg.dbQueries.GetAllQuotes(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var tmplFile = "templates/allquotes.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = tmpl.Execute(w, records)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
