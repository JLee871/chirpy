package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonResp(w http.ResponseWriter, code int, jStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(jStruct)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

func errorResp(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	errResp := errorResponse{
		Error: msg,
	}

	jsonResp(w, code, errResp)
}
