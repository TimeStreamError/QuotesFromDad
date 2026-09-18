package main

import (
	"html/template"
	"net/http"
)

func (apiCfg *APIConfig) handleGetRandom(w http.ResponseWriter, r *http.Request) {
	randomQuote, err := apiCfg.dbQueries.GetRandom(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	author := ""
	if randomQuote.Author.Valid {
		author = randomQuote.Author.String
	}
	type TemplateData struct {
		Quote  string
		Author string
	}
	data := TemplateData{}
	data.Author = author
	data.Quote = randomQuote.Quote

	var tmplFile = "templates/randomquote.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
