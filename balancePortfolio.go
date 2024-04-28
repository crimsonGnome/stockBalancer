package main

import (
	"encoding/json"
	"net/http"
)

func BalancePortfolio(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Balance Portfolio feature called"})
}
// 