package main

import (
	"fmt"
	AthenaCall "my-lambda-app/athena"
	dynamo "my-lambda-app/dynamo"

	"github.com/aws/aws-lambda-go/events"
)

func CalculateBackTest() events.APIGatewayProxyResponse {
	client := dynamo.GetClient()
	date := "2019-11-14"
	currentStockPrices := AthenaCall.SQL_date_Price(date)

	currentCash := dynamo.GetCashTotal("historicalCashBackTest", client)
	currentPortfolio := dynamo.GetPortfolio("historicalTestPortfolio", client)
	totalCash := currentCash.Amount + CalculateCashTotal(currentPortfolio, currentStockPrices)

	fmt.Println("Total Cash:", totalCash)
	message := fmt.Sprintf("Total cash %f", totalCash)
	// Example processing for the /backTest path
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonMessage(message),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
