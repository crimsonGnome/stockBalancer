package dynamo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Orders struct {
	Symbol    string
	Date      string
	OrderType string
	Amount    int
}

type PredictionData struct {
	Symbol         string
	PredictedPrice float64
	Bias           float64
	MedianDistance float64
	OriginalPrice  float64
	GrowthRate     float64
}

func BatchWriteTransactionHistory(portfolioInput []Orders, client *dynamodb.Client) error {
	writeRequests := make([]types.WriteRequest, len(portfolioInput))

	for i, order := range portfolioInput {
		item := map[string]types.AttributeValue{
			"date":      &types.AttributeValueMemberS{Value: order.Date},
			"symbol":    &types.AttributeValueMemberS{Value: order.Symbol},
			"orderType": &types.AttributeValueMemberS{Value: order.OrderType},
			"amount":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", order.Amount)},
		}
		writeRequests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		}
	}

	batchInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"stock-transaction-history": writeRequests,
		},
	}

	result, err := client.BatchWriteItem(context.TODO(), batchInput)
	if err != nil {
		log.Fatalf("failed to batch write items: %v", err)
		return err
	}
	fmt.Printf("Batch write successful: %v\n", result)

	return nil
}

func BatchWritePredictionWeights(predictionWeight []PredictionData, client *dynamodb.Client) error {
	writeRequests := make([]types.WriteRequest, 15)
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02")

	for i := 0; i < 15; i++ {
		item := map[string]types.AttributeValue{
			"date":           &types.AttributeValueMemberS{Value: formattedTime},
			"symbol":         &types.AttributeValueMemberS{Value: predictionWeight[i].Symbol},
			"weight":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", predictionWeight[i].GrowthRate)},
			"predictedPrice": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", predictionWeight[i].PredictedPrice)},
			"price":          &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", predictionWeight[i].OriginalPrice)},
			"bias":           &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", predictionWeight[i].Bias)},
			"medianDistance": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", predictionWeight[i].MedianDistance)},
			"rank":           &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i)},
		}
		writeRequests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: item,
			},
		}
	}

	batchInput := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			"prediction-weights": writeRequests,
		},
	}

	result, err := client.BatchWriteItem(context.TODO(), batchInput)
	if err != nil {
		log.Fatalf("failed to batch write items: %v", err)
		return err
	}
	fmt.Printf("Batch write successful: %v\n", result)

	return nil

}
