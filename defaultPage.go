package main

import (
	"github.com/aws/aws-lambda-go/events"
)

func HandleDefault() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       jsonMessage("Welcome to the API"),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
