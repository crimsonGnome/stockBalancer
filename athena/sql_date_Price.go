package Athena

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/athena"
    "github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func SQL_date_Price(date srting)  {
    // Load the AWS SDK for Go configuration
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatalf("unable to load SDK config, %v", err)
    }

    // Create an Athena client with the configuration
    client := athena.NewFromConfig(cfg)

    // SQL query to execute
    query := fmt.sprinf("SELECT * FROM updated_csv WHERE datetime =  '%s';", date)
    database := "combined-stock-data-crimson" // Specify your Athena database
    resultConfiguration := types.ResultConfiguration{
        OutputLocation: aws.String("s3://stock-trading-bucket-crimson/"), // Specify your S3 output path
    }

    // Execute the query
    input := &athena.StartQueryExecutionInput{
        QueryString:         aws.String(query),
        QueryExecutionContext: &types.QueryExecutionContext{
            Database: aws.String(database),
        },
        ResultConfiguration: &resultConfiguration,
    }

    result, err := client.StartQueryExecution(context.TODO(), input)
    if err != nil {
        log.Fatalf("failed to execute query, %v", err)
    }

    fmt.Printf("Query Execution ID: %s\n", *result.QueryExecutionId)

    return result 
}