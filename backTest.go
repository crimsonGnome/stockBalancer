package main

import (
	"fmt"

	AthenaCall "my-lambda-app/athena"
	dynamo "my-lambda-app/dynamo"
	helperFunctions "my-lambda-app/helperFunctions"

	"github.com/aws/aws-lambda-go/events"
)

func HandleBackTest() events.APIGatewayProxyResponse {
	// Loop through back data Data
	// Define back test Data Loop
	date := "2017-02-14"

	// Loop over data ( build first with just one data set)
	// call the SQL athena function to get actual values
	currentStockPrices := AthenaCall.SQL_date_Price(date)
	// ------------------ Get Current cash Balance ----------------------------------------------

	// Call Get Current Portfolio distribution
	currentPortfolio := dynamo.GetPortfolio("historicalTestPortfolio")
	totalCash := 0.0
	fmt.Println("Printing current portfolio:", *currentPortfolio)
	if len(*currentPortfolio) == 0 {
		// Get current cash by accessing cashTable
		currentCash := dynamo.GetCashTotal("historicalCashBackTest")
		totalCash = currentCash.Amount

		fmt.Println("Printing current cash:", currentCash)
		fmt.Println("Total Cash:", totalCash)
	} else {
		currentCash := dynamo.GetCashTotal("historicalCashBackTest")
		totalCash = currentCash.Amount + CalculateCashTotal(currentPortfolio, currentStockPrices)
		fmt.Println("Total Cash:", totalCash)
	}

	// --------------------------------------------------------------------------------

	// Get Predictions  in the range of Dynamo table (Eliminate all stocks that are below .5 bias)

	// Call tradingStrategy take in all the data to get (currentCash, prediction, current prices)
	// 18 + 16 + 14 + 12 + 10 + 8 + 6 + 6 + 6 + 4

	// rectify_portfolio - buy and sell data into new table - save logs in portfolio transaction

	// save new Portfolio distribution in dynamo table

	// Update left over cash

	// Example processing for the /backTest path
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonMessage("BackTest process started"),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func CalculateCashTotal(currentPortfolio *[]dynamo.PortfolioDistribution, currentStockPrices *[]map[string]string) float64 {
	totalCash := 0.0
	portfolioMap := make(map[string]float64)
	for i := 0; i < len(*currentPortfolio); i++ {
		portfolioMap[(*currentPortfolio)[i].Symbol] = (*currentPortfolio)[i].Amount
	}

	for _, record := range *currentStockPrices {
		amount, ok := portfolioMap[record["symbol"]]
		if !ok {
			continue
		}

		fmt.Println("In portfolio:", record["symbol"])
		totalCash = totalCash + (amount * helperFunctions.FloatConverter(record["close"]))
	}

	return totalCash

}
