package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Detail struct {
	UserID string `json:"UserID"`
}

type Event struct {
	Source     string `json:"source"`
	Detail     Detail `json:"detail"`
	DetailType string `json:"detailType"`
}

func HandleRequest(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("received nil event")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := dynamodb.NewFromConfig(cfg)

	for {
		// Initialize the keysToDelete slice to hold keys for batch deletion
		keysToDelete := make([]map[string]types.AttributeValue, 0)

		// Define the Scan input with the filter expression
		input := &dynamodb.ScanInput{
			TableName:        aws.String("KwetterMessagesDynamoDBTable"),
			FilterExpression: aws.String("UserID = :userID"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":userID": &types.AttributeValueMemberS{Value: event.Detail.UserID},
			},
		}

		// Paginate through the scan results and collect keys for batch deletion
		paginator := dynamodb.NewScanPaginator(client, input)
		for paginator.HasMorePages() {
			page, err := paginator.NextPage(ctx)
			if err != nil {
				return err
			}

			for _, item := range page.Items {
				keysToDelete = append(keysToDelete, map[string]types.AttributeValue{
					"UserID": item["UserID"],
					"GUID":   item["GUID"],
				})
			}
		}

		// If no items are left, break out of the loop
		if len(keysToDelete) == 0 {
			break
		}

		// Batch delete items in groups of the maximum allowed limit
		for len(keysToDelete) > 0 {
			// Slice the next batch of keys to delete
			batchSize := 25 // Maximum batch size
			if len(keysToDelete) < batchSize {
				batchSize = len(keysToDelete)
			}
			keys := keysToDelete[:batchSize]
			keysToDelete = keysToDelete[batchSize:]

			// Prepare the batch delete input
			requestItems := map[string][]types.WriteRequest{
				"KwetterMessagesDynamoDBTable": make([]types.WriteRequest, len(keys)),
			}
			for i, key := range keys {
				requestItems["KwetterMessagesDynamoDBTable"][i] = types.WriteRequest{
					DeleteRequest: &types.DeleteRequest{Key: key},
				}
			}

			// Execute the batch delete request with retry logic
			for attempts := 0; attempts < 3; attempts++ {
				_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
					RequestItems: requestItems,
				})
				if err != nil {
					// If the error is a throttling exception, retry after waiting for some time
					if isThrottleException(err) {
						fmt.Println("Request throttled. Retrying...")
						time.Sleep(1 * time.Second) // Wait for 1 second before retrying
						continue
					}
					return err
				}
				// Batch deletion successful, break out of retry loop
				break
			}
		}
	}

	return nil
}

// isThrottleException checks if the error message indicates a throttling exception
func isThrottleException(err error) bool {
	// Check the error message for indications of throttling
	return err != nil && (contains(err.Error(), "RequestLimitExceeded") || contains(err.Error(), "ThrottlingException"))
}

// contains checks if a string contains another substring
func contains(s, substr string) bool {
	return s != "" && substr != "" && len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}

func main() {
	lambda.Start(HandleRequest)
}
