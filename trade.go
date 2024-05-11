package main

import (
	"github.com/aws/aws-lambda-go/events"
)

func HandleTrade() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonMessage("Trade executed"),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

//
