package dynamo

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PortfolioDistribution struct {
	Symbol string
	Amount float64
}

// Initialize a DynamoDB client
func GetClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	return dynamodb.NewFromConfig(cfg)
}

// Fetch all items from the specified table
func GetPortfolio(recordType string, client *dynamodb.Client) *[]PortfolioDistribution {
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

	var distributions []PortfolioDistribution

	val := result.Item["portfolio"].(*types.AttributeValueMemberSS)
	fmt.Printf("Retrieved item: %v\n", result.Item)
	if val.Value[0] == "empty" {
		fmt.Println("hit- empty")
		return &distributions
	}
	fmt.Println("val from portfolio:", val)
	for i := 0; i < len(val.Value); i++ {
		var tempDistribution PortfolioDistribution
		parts := strings.Split(val.Value[i], ":")
		symbol := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		tempDistribution.Symbol = symbol
		tempDistribution.Amount, err = strconv.ParseFloat(value, 64)
		if err != nil {
			log.Fatalf("Failed to parse Amount, %v", err)
		}
		fmt.Println("Stock distribution:", val.Value[i])
		distributions = append(distributions, tempDistribution)
	}
	return &distributions
}
