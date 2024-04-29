package dynamo

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetCashTotal(recordType string) PortfolioDistribution {
	client := getClient()
	// Specify the table name and the primary key of the item
	tableName := "stock-current-portfolio"
	itemKey := map[string]types.AttributeValue{
		"RecordType": &types.AttributeValueMemberS{Value: recordType},
	}

	// Create the GetItem input
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       itemKey,
	}

	// Call DynamoDB to get the item
	result, err := client.GetItem(context.Background(), input)
	if err != nil {
		log.Fatalf("failed to get item from DynamoDB, %v", err)
	}

	// Check if the result is empty
	if result.Item == nil {
		log.Println("No item found with the specified key")
	} else {
		// Print the retrieved item
		fmt.Printf("Retrieved item: %v\n", result.Item)
	}
	var distribution PortfolioDistribution

	distribution.Symbol = recordType
	val := result.Item["Amount"].(*types.AttributeValueMemberN)
	distribution.Amount, err = strconv.ParseFloat(val.Value, 64)
	if err != nil {
		log.Fatalf("Failed to parse Amount, %v", err)
	}

	return distribution
}
