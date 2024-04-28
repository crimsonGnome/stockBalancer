package main

import (
	"encoding/json"
	"net/http"
)

func DefaultPage(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to the default page"})
}
