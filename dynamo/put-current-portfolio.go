package dynamo

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func PutCurrentPortfolio(key string, portfolioInput []string, client *dynamodb.Client) error {
	item := map[string]types.AttributeValue{
		"RecordType": &types.AttributeValueMemberS{Value: key}, // Assuming 'ID' is the primary key
		// Add other attributes here
		"portfolio": &types.AttributeValueMemberSS{Value: portfolioInput},
	}

	// Create the PutItem input
	input := &dynamodb.PutItemInput{
		TableName: aws.String("stock-current-portfolio"),
		Item:      item,
	}

	// Execute the PutItem operation
	_, err := client.PutItem(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to put item: %v", err)
		return err
	}

	fmt.Println("Successfully put item in table")

	return nil

}

func PutCurrentCash(key string, cashUpdate float64, client *dynamodb.Client) error {
	item := map[string]types.AttributeValue{
		"RecordType": &types.AttributeValueMemberS{Value: key}, // Assuming 'ID' is the primary key
		// Add other attributes here
		"Amount": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", cashUpdate)},
	}

	// Create the PutItem input
	input := &dynamodb.PutItemInput{
		TableName: aws.String("stock-current-portfolio"),
		Item:      item,
	}

	// Execute the PutItem operation
	_, err := client.PutItem(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to put item: %v", err)
		return err
	}

	fmt.Println("Successfully put item in table")

	return nil

}
