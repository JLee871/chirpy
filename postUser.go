// Handler for creating new users
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JLee871/chirpy/internal/auth"
	"github.com/JLee871/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (c *apiConfig) postuserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResp(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbUser, err := c.databaseQueries.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashedPW})
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Could not create user", err)
		return
	}

	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	jsonResp(w, http.StatusCreated, response{newUser})
}
