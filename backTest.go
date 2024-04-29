package main

import (
	"fmt"
	AthenaCall "my-lambda-app/athena"

	"github.com/aws/aws-lambda-go/events"
)

func HandleBackTest() events.APIGatewayProxyResponse {
	// Loop through back data Data
	// Define back test Data Loop
	date := "2017-02-14"

	// Loop over data ( build first with just one data set)
	// call the SQL athena function to get actual values
	currentStockPrices := AthenaCall.SQL_date_Price(date)
	for _, result := range currentStockPrices {
		fmt.Println(result)
	}
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

	// Example processing for the /backTest path
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonMessage("BackTest process started"),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
