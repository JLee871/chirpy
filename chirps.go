// Handlers for CRUD chirps operations
package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/JLee871/chirpy/internal/auth"
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

func (c *apiConfig) postchirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	//Authentication steps
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "bad request", err)
		return
	}
	userID, err := auth.ValidateJWT(token, c.tokenSecret)
	if err != nil {
		errorResp(w, http.StatusUnauthorized, "status unauthorized", errors.New("token validation failed"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Couldn't decode parameters", err)
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
		UserID: userID,
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

func (c *apiConfig) getallchirpsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("author_id")
	sortType := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error
	var ID uuid.UUID
	if userID == "" {
		chirps, err = c.databaseQueries.GetChirps(r.Context())
	} else {
		ID, err = uuid.Parse(userID)
		if err != nil {
			errorResp(w, http.StatusBadRequest, "could not parse author id", err)
			return
		}
		chirps, err = c.databaseQueries.GetChirpsByUser(r.Context(), ID)
	}

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

	if sortType == "desc" {
		sort.Slice(resp, func(i, j int) bool {
			return resp[i].CreatedAt.After(resp[j].CreatedAt)
		})
	}

	jsonResp(w, http.StatusOK, resp)
}

func (c *apiConfig) getsinglechirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Could not parse chirp ID", err)
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

func (c *apiConfig) delchirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResp(w, http.StatusUnauthorized, "could not find token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, c.tokenSecret)
	if err != nil {
		errorResp(w, http.StatusUnauthorized, "could not validate token", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Could not parse chirp ID", err)
		return
	}

	dbChirp, err := c.databaseQueries.GetChirp(r.Context(), chirpID)
	if err != nil {
		errorResp(w, http.StatusNotFound, "Could not find chirp", err)
		return
	}

	if dbChirp.UserID != userID {
		errorResp(w, http.StatusForbidden, "user is not author of chirp", err)
		return
	}

	err = c.databaseQueries.DeleteChirp(r.Context(), database.DeleteChirpParams{ID: chirpID, UserID: userID})
	if err != nil {
		errorResp(w, http.StatusBadRequest, "could not delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
