package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/timestreamerror/capstone/internal/auth"
	"github.com/timestreamerror/capstone/internal/database"
)

func (apiCfg *APIConfig) handleSignup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type ParsedData struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		Authorization string `json:"authorization"`
	}

	data := ParsedData{}
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if data.Email == "" || data.Password == "" || data.Authorization == "" {
		respondWithError(w, http.StatusBadRequest, "Email, password and authorization required")
		return
	}

	if data.Authorization != apiCfg.loginAuthorization {
		respondWithError(w, http.StatusBadRequest, "Incorrect authorization")
		return
	}

	params := database.CreateUserParams{}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)

	params.ID = newUUID
	params.CreatedAt = time.Now()
	params.UpdatedAt = time.Now()
	params.Email = data.Email
	params.HashedPassword = hashedPassword

	_, err = apiCfg.dbQueries.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
