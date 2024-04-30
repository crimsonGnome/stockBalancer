package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	switch request.RawPath {
	case "/":
		return HandleDefault(), nil
	case "/backTest":
		return HandleBackTest(), nil
	case "/trade":
		return HandleTrade(), nil
	case "/calculate":
		return CalculateBackTest(), nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}, nil
	}
}

func jsonMessage(message string) string {
	response := Response{Message: message}
	responseBody, err := json.Marshal(response)
	if err != nil {
		return "{\"message\":\"Error creating response\"}"
	}
	return string(responseBody)
}

func main() {
	lambda.Start(Handler)
}
