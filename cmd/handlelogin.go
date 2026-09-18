package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/timestreamerror/capstone/internal/auth"
	"github.com/timestreamerror/capstone/internal/database"
)

func (apiCfg *APIConfig) handlerLogin(res http.ResponseWriter, req *http.Request) {
	type Params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type UsersReturn struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	data := Params{}
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	decodedEmail := data.Email
	decodedPassword := data.Password

	queriedUser, err := apiCfg.dbQueries.LoginViaEmail(req.Context(), decodedEmail)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	passwordMatches, err := auth.CheckPasswordHash(decodedPassword, queriedUser.HashedPassword)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	if !passwordMatches {
		respondWithError(res, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	returnUser := UsersReturn{
		ID:        queriedUser.ID,
		CreatedAt: queriedUser.CreatedAt,
		UpdatedAt: queriedUser.UpdatedAt,
		Email:     queriedUser.Email,
	}

	token, err := auth.MakeJWT(queriedUser.ID, apiCfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	returnUser.Token = token
	refreshToken := auth.MakeRefreshToken()
	refreshTokenParams := database.InsertRefreshTokenParams{}
	refreshTokenParams.Token = refreshToken
	refreshTokenParams.UserID = queriedUser.ID
	refreshTokenParams.ExpiresAt = time.Now().Add(time.Hour * 24 * 60)
	dbRefreshToken, err := apiCfg.dbQueries.InsertRefreshToken(req.Context(), refreshTokenParams)
	if err != nil {
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	returnUser.RefreshToken = dbRefreshToken.Token

	http.SetCookie(res, &http.Cookie{
		Name:     "refreshtoken",
		Value:    dbRefreshToken.Token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30, // 30 days
	})

	respondWithJSON(res, http.StatusOK, returnUser)
}
