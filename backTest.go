package main

import (
	"encoding/json"
	"net/http"
)

func BackTest(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "BackTest feature called"})
}
