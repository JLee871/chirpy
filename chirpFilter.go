package main

import "strings"

var badWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

func filterMsg(msg string) string {
	msgWords := strings.Split(msg, " ")
	for i, word := range msgWords {
		if badWords[strings.ToLower(word)] {
			msgWords[i] = "****"
		}
	}
	return strings.Join(msgWords, " ")
}
