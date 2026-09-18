package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/timestreamerror/capstone/internal/database"
)

func (apiCfg *APIConfig) handleAddTag(w http.ResponseWriter, r *http.Request) {

	// adds a tag to the tags table
	params := database.AddTagParams{}
	id, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	params.ID = id
	params.CreatedAt = time.Now()
	params.UpdatedAt = time.Now()
	decoder := json.NewDecoder(r.Body)

	type DataJSON struct {
		Name string `json:"name"`
	}
	data := DataJSON{}

	err = decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	params.Name = data.Name
	tagData, err := apiCfg.dbQueries.AddTag(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tagData)
}
