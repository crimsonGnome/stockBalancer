package dynamo

import (
	"context"
	helperFunctions "my-lambda-app/helperFunctions"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PredictionPriceQuery struct {
	Symbol         string
	Price          float64
	Bias           float64
	MedianDistance float64
	MeanDistance   float64
	DefaultError   float64
}

func ScanPredictions(dateInput string, client *dynamodb.Client, ctx context.Context) (*[]PredictionPriceQuery, error) {
	tableName := "stock-predictions"
	var predictedPrices []PredictionPriceQuery
	input := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("#date = :date"),
		ExpressionAttributeNames: map[string]string{
			"#date": "date", // Alias for the attribute name to avoid conflict with reserved words
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":date": &types.AttributeValueMemberS{Value: dateInput},
		},
	}
	// Execute the Scan
	resp, err := client.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	for _, itemMap := range resp.Items {
		var temp PredictionPriceQuery
		// Extract and convert other attributes manually
		if val, ok := itemMap["symbol"].(*types.AttributeValueMemberS); ok {
			temp.Symbol = val.Value
		}
		if val, ok := itemMap["price"].(*types.AttributeValueMemberN); ok {
			temp.Price = helperFunctions.FloatConverter(val.Value)

		}
		if val, ok := itemMap["bias"].(*types.AttributeValueMemberN); ok {
			temp.Bias = helperFunctions.FloatConverter(val.Value)
		}
		if val, ok := itemMap["medianDistance"].(*types.AttributeValueMemberN); ok {
			temp.MedianDistance = helperFunctions.FloatConverter(val.Value)
		}
		if val, ok := itemMap["meanDistance"].(*types.AttributeValueMemberN); ok {
			temp.MeanDistance = helperFunctions.FloatConverter(val.Value)
		}
		if val, ok := itemMap["defaultError"].(*types.AttributeValueMemberN); ok {
			temp.DefaultError = helperFunctions.FloatConverter(val.Value)
		}

		predictedPrices = append(predictedPrices, temp)
	}

	return &predictedPrices, nil

}
