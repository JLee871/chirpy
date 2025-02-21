// Handler for creating new chirps
package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/JLee871/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (c *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//Check if chirp exceeds max length
	const maxChirpLen = 140
	if len(params.Body) > maxChirpLen {
		errorResp(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	dbChirp, err := c.databaseQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   filterMsg(params.Body),
		UserID: params.UserID,
	})

	if err != nil {
		errorResp(w, http.StatusBadRequest, "Could not create chirp", err)
		return
	}

	newChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	jsonResp(w, http.StatusCreated, response{newChirp})
}
