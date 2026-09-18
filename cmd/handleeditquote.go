package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/timestreamerror/capstone/internal/database"
)

func (apiCfg *APIConfig) handleEditQuote(w http.ResponseWriter, r *http.Request) {
	quoteIDString := chi.URLParam(r, "quote-id")
	quoteID, err := uuid.Parse(quoteIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	quote, err := apiCfg.dbQueries.EditQuote(r.Context(), quoteID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var tmplFile = "templates/editquote.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type Tag struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	type TemplateParams struct {
		Quote  string                   `json:"quote"`
		Author string                   `json:"author"`
		ID     string                   `json:"id"`
		Tags   []database.GetAllTagsRow `json:"tags"`
	}
	templateParams := TemplateParams{}

	tagNames, err := apiCfg.dbQueries.GetAllTags(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	templateParams.Tags = tagNames
	templateParams.ID = quoteIDString
	author := ""
	if quote.Author.Valid {
		author = quote.Author.String
	}
	templateParams.Author = author
	templateParams.Quote = quote.Quote

	err = tmpl.Execute(w, templateParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
