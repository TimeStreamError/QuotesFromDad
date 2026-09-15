package main

import (
	"net/http"
	"time"

	"github.com/timestreamerror/capstone/internal/auth"
)

func (apiCfg *APIConfig) handleRefreshTokenCheck(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshtoken")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	dbRefreshToken, err := apiCfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken.Value)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	if !dbRefreshToken.ExpiresAt.After(time.Now()) || dbRefreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Error: refresh token expired or revoked")
		return
	}

	token, err := auth.MakeJWT(dbRefreshToken.UserID, apiCfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	type DataToken struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, DataToken{Token: token})
}
