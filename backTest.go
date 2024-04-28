package main

import (
	"encoding/json"
	"net/http"
	"my-lambda-app/athena"
)

func BackTest(w http.ResponseWriter, r *http.Request) {
	// Loop through back data Data 
	// Define back test Data Loop
	date := "2017-02-14"

	// Loop over data ( build first with just one data set) 
	// call the SQL athena function to get actual values
	currentStockPrices := Athena.SQL_date_Price(date)

	// ------------------ Get Current cash Balance ----------------------------------------------

	// Call Get Current Portolio distribution 
	// calulate Current Portoflio Prices based on athena data and my distribution ()
	// Get current cash by adding leftOver cash table + current Portfolio market rate
	
	// --------------------------------------------------------------------------------

	// Get Predictions  in the range of Dynamo table (Elimnate all stocks that are below .5 bias)

	// Call tradingStrategy take in all the data to get (currentCash, prediction, current prices)

	// rectify_portfolio - buy and sell data into new table - save logs in portfolio transaction

	// save new portoflio distrubution in dynamo table 
	
	// Update left over cash









	json.NewEncoder(w).Encode(map[string]string{"message": currentStockPrices})
}
