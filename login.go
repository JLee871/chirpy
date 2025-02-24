package main

import (
	"encoding/json"
	"net/http"

	"github.com/JLee871/chirpy/internal/auth"
)

func (c *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
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

	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	jsonResp(w, http.StatusOK, newUser)
}
