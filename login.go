package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JLee871/chirpy/internal/auth"
	"github.com/JLee871/chirpy/internal/database"
)

// Default token expiration time of 1 hour (in seconds)
const defaultExpiration = 3600

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type loginResp struct {
		User
		Token         string `json:"token"`
		Refresh_token string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	dbUser, err := c.databaseQueries.GetUserFromEmail(r.Context(), params.Email)
	if err != nil {
		errorResp(w, http.StatusNotFound, "Couldn't find email", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		errorResp(w, http.StatusUnauthorized, "wrong password", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, c.tokenSecret, time.Duration(defaultExpiration)*time.Second)
	if err != nil {
		errorResp(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		errorResp(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	dbToken, err := c.databaseQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: dbUser.ID})
	if err != nil {
		errorResp(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	newResp := loginResp{
		User: User{
			ID:          dbUser.ID,
			CreatedAt:   dbUser.CreatedAt,
			UpdatedAt:   dbUser.UpdatedAt,
			Email:       dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token:         accessToken,
		Refresh_token: dbToken.Token,
	}

	jsonResp(w, http.StatusOK, newResp)
}

func (c *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type refreshResp struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResp(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	dbToken, err := c.databaseQueries.GetToken(r.Context(), refreshToken)

	if err != nil || time.Now().After(dbToken.ExpiresAt) || dbToken.RevokedAt.Valid {
		errorResp(w, http.StatusUnauthorized, "status unauthorized", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, c.tokenSecret, time.Duration(defaultExpiration)*time.Second)
	if err != nil {
		errorResp(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	jsonResp(w, http.StatusOK, refreshResp{Token: accessToken})
}

func (c *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "could not find token", err)
		return
	}

	err = c.databaseQueries.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "could not revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
