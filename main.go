package main

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
)

// Handler is your Lambda function handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize the mux router
	r := mux.NewRouter()
	r.HandleFunc("/", DefaultPage).Methods("GET")
	r.HandleFunc("/balancePortfolio", BalancePortfolio).Methods("POST")
	r.HandleFunc("/backTest", BackTest).Methods("POST")

	// The incoming request is adapted to an http.Request
	httpRequest, err := http.NewRequest(req.HTTPMethod, req.Path, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Using a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httpRequest)

	// Convert the response recorder's result into APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{
		StatusCode: rr.Code,
		Body:       rr.Body.String(),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func main() {
	// Start the Lambda handler
	lambda.Start(Handler)
}
