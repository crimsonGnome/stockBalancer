package main

import (
	"context"
	"log"
	AthenaCall "my-lambda-app/athena"
	dynamo "my-lambda-app/dynamo"
	"my-lambda-app/env"

	"github.com/aws/aws-lambda-go/events"
)

func HandleDefault() events.APIGatewayProxyResponse {
	client := dynamo.GetClient()
	currentPredictionFlag := true
	dateStringDaily := env.ENV_CURRENT_DATE
	predictionArray := env.ENV_FUTURE_PREDICTIONS

	// // Convert current date TIme string into a Date time
	// dateTimeDaily, err := time.Parse("2006-01-02", dateStringDaily)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Compare times to adjust
	currentStockPrices := AthenaCall.SQL_date_Price(dateStringDaily)

	// Call Get Current Portfolio distribution
	currentPortfolio := dynamo.GetPortfolio("currentPortfolio", client)
	totalCash := 0.0

	// Is there a currentPortfolio
	if len(*currentPortfolio) == 0 {
		// Get current cash by accessing cashTable
		currentCash := dynamo.GetCashTotal("currentCash", client)
		totalCash = currentCash.Amount
	} else {
		currentCash := dynamo.GetCashTotal("currentCash", client)
		totalCash = currentCash.Amount + CalculateCashTotal(currentPortfolio, currentStockPrices)
	}

	// Get Predictions  in the range of Dynamo table (Eliminate all stocks that are below .5 bias)
	predictions := TradingWeight(*currentStockPrices, predictionArray, 0, client, context.TODO())

	// Call UpdatePortfolio - creates a new portfolio
	updatedPortfolio, cashUpdate := UpdatePortfolio(totalCash, &predictions, currentPredictionFlag, client)

	// rectify_portfolio - buy and sell data into new table - save logs in portfolio transaction
	buys, sells := RectifyPortfolio(dateStringDaily, currentPortfolio, &updatedPortfolio)

	// save new Portfolio distribution in dynamo table
	stringPortfolio := ConvertPortToString(updatedPortfolio)
	err := dynamo.PutCurrentPortfolio("currentPortfolio", stringPortfolio, client)
	if err != nil {
		log.Fatalf("Failed to PutCurrentPortfolio: %v", err)
	}

	// Update Transaction Trade Logs
	// Update Cash Table
	dynamo.PutCurrentCash("currentCash", cashUpdate, client)
	if err != nil {
		log.Fatalf("BatchWriteTransactionHistory Sells: %v", err)
	}

	// Update Sells
	if len(*sells) != 0 {
		dynamo.BatchWriteTransactionHistory(*sells, client)
		if err != nil {
			log.Fatalf("BatchWriteTransactionHistory Sells: %v", err)
		}
	}

	// Update Buys
	dynamo.BatchWriteTransactionHistory(*buys, client)
	if err != nil {
		log.Fatalf("BatchWriteTransactionHistory Sells: %v", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonMessage("Balanced Portfolio"),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
