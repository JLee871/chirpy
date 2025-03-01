package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/JLee871/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) polkawebhookHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		errorResp(w, http.StatusUnauthorized, "could not retrieve api key", err)
		return
	}

	if apiKey != os.Getenv("POLKA_KEY") {
		errorResp(w, http.StatusUnauthorized, "bad api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "could not parse user id", err)
		return
	}

	_, err = c.databaseQueries.GetUserFromID(r.Context(), userID)
	if err != nil {
		errorResp(w, http.StatusNotFound, "could not find user", err)
		return
	}

	err = c.databaseQueries.UpgradeUserRed(r.Context(), userID)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "could not upgrade user to red", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
