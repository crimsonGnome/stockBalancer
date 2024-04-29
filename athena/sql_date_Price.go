package AthenaCall

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func SQL_date_Price(date string) *[]map[string]string {
	ctx := context.TODO()

	// Load the AWS SDK for Go configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return nil
	}

	// Create an Athena client with the configuration
	client := athena.NewFromConfig(cfg)

	// SQL query to execute
	query := fmt.Sprintf("SELECT * FROM updated_csv WHERE date_parse(datetime, '%%Y-%%m-%%d') = date '%s'", date)
	// TEST SQL - query := fmt.Sprintf("SELECT * FROM `combined-stock-data-crimson`.`updated_csv` LIMIT 10")
	database := "combined-stock-data-crimson"
	resultConfiguration := types.ResultConfiguration{
		OutputLocation: aws.String("s3://stock-trading-bucket-crimson/"),
	}

	// Execute the query
	input := &athena.StartQueryExecutionInput{
		QueryString: aws.String(query),
		QueryExecutionContext: &types.QueryExecutionContext{
			Database: aws.String(database),
		},
		ResultConfiguration: &resultConfiguration,
	}

	startResult, err := client.StartQueryExecution(ctx, input)
	if err != nil {
		log.Fatalf("failed to execute query, %v", err)
		return nil
	}

	queryExecutionID := startResult.QueryExecutionId

	// Wait for the query to complete
	for {
		execInput := &athena.GetQueryExecutionInput{
			QueryExecutionId: queryExecutionID,
		}
		execOutput, err := client.GetQueryExecution(ctx, execInput)
		if err != nil {
			log.Fatalf("failed to get query execution, %v", err)
			return nil
		}

		status := execOutput.QueryExecution.Status.State
		if status == types.QueryExecutionStateSucceeded {
			break
		} else if status == types.QueryExecutionStateFailed || status == types.QueryExecutionStateCancelled {
			log.Fatalf("query failed or was cancelled: %v", *execOutput.QueryExecution.Status.StateChangeReason)
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	// Fetch the results
	results := []map[string]string{}
	resultInput := &athena.GetQueryResultsInput{
		QueryExecutionId: queryExecutionID,
	}

	paginator := athena.NewGetQueryResultsPaginator(client, resultInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to get query results page, %v", err)
			return nil
		}

		for _, row := range page.ResultSet.Rows {
			if len(page.ResultSet.ResultSetMetadata.ColumnInfo) != len(row.Data) {
				continue // Skip header or mismatched rows
			}
			resultMap := make(map[string]string)
			for idx, col := range row.Data {
				columnName := *page.ResultSet.ResultSetMetadata.ColumnInfo[idx].Name // Dereference here
				if col.VarCharValue != nil {
					resultMap[columnName] = *col.VarCharValue
				}
			}
			results = append(results, resultMap)
		}

	}

	return &results
}
