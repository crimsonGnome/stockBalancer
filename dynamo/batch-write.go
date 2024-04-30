package dynamo

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Orders struct {
	Symbol    string
	Date      string
	OrderType string
	Amount    int
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
