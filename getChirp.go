package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (c *apiConfig) getallchirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.databaseQueries.GetChirps(r.Context())
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "Could not retrieve chirps", err)
		return
	}
	resp := []Chirp{}

	for _, dbChirp := range chirps {
		resp = append(resp, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	jsonResp(w, http.StatusOK, resp)
}

func (c *apiConfig) getsinglechirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "Could not parse chirp ID", err)
		return
	}

	dbChirp, err := c.databaseQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		errorResp(w, http.StatusNotFound, "Could not find chirp", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	jsonResp(w, http.StatusOK, chirp)
}
