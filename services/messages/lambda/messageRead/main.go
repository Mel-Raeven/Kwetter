package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

type Message struct {
	GUID    string `json:"GUID"`
	Message string `json:"Message"`
	Ts      string `json:"ts"`
	Userid  string `json:"UserID"`
}

func getMessageHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["userid"]

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// Create DynamoDB client
	client := dynamodb.NewFromConfig(cfg)

	// Extract last evaluated key from request headers
	lastEvaluatedKey := request.Headers["last-evaluated-key"]
	fmt.Print(lastEvaluatedKey)
	// Create key condition expression for the query
	keyCondExp := "UserID = :userid"
	expAttrValues := map[string]types.AttributeValue{
		":userid": &types.AttributeValueMemberS{Value: id},
	}

	// Create query input with key condition expression, tablename, and projection expression
	input := &dynamodb.QueryInput{
		TableName:                 aws.String("KwetterMessagesDynamoDBTable"),
		KeyConditionExpression:    aws.String(keyCondExp),
		ExpressionAttributeValues: expAttrValues,
		ProjectionExpression:      aws.String("GUID, Message, ts, UserID"),
		Limit:                     aws.Int32(10), // Limit to 10 items
	}

	// Set exclusive start key if available
	if lastEvaluatedKey != "" {
		var lastEvaluatedMap map[string]map[string]string
		if err := json.Unmarshal([]byte(lastEvaluatedKey), &lastEvaluatedMap); err != nil {
			fmt.Printf("Failed to unmarshal last evaluated key: %v\n", err)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, err
		}

		lastEvaluatedDynamoMap := make(map[string]types.AttributeValue)
		for k, v := range lastEvaluatedMap {
			switch {
			case v["S"] != "":
				lastEvaluatedDynamoMap[k] = &types.AttributeValueMemberS{Value: v["S"]}
			case v["N"] != "":
				lastEvaluatedDynamoMap[k] = &types.AttributeValueMemberN{Value: v["N"]}
			case v["B"] != "":
				lastEvaluatedDynamoMap[k] = &types.AttributeValueMemberB{Value: []byte(v["B"])}
			}
		}

		input.ExclusiveStartKey = lastEvaluatedDynamoMap
	}

	// Perform the query
	result, err := client.Query(ctx, input)
	if err != nil {
		fmt.Printf("Couldn't query info for user %v. Here's why: %v\n", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	if len(result.Items) == 0 {
		fmt.Printf("No messages found for user %v\n", id)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "*",
				"Access-Control-Allow-Headers": "*",
				"Content-Type":                 "application/json",
			},
			Body: "[]", // Return empty array if no messages found
		}, nil
	}

	// Create list of messages from results
	var messages []Message
	err = attributevalue.UnmarshalListOfMaps(result.Items, &messages)
	if err != nil {
		fmt.Printf("Couldn't unmarshal query response. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// Prepare the last evaluated key for the response header
	lastEvaluatedKeyString := ""
	if result.LastEvaluatedKey != nil {
		// Create a map with the correct attribute types
		correctedLastEvaluatedKey := make(map[string]map[string]string)
		for k, v := range result.LastEvaluatedKey {
			switch v := v.(type) {
			case *types.AttributeValueMemberS:
				correctedLastEvaluatedKey[k] = map[string]string{"S": v.Value}
			case *types.AttributeValueMemberN:
				correctedLastEvaluatedKey[k] = map[string]string{"N": v.Value}
			case *types.AttributeValueMemberB:
				correctedLastEvaluatedKey[k] = map[string]string{"B": string(v.Value)}
			}
		}
		marshalledKey, err := json.Marshal(correctedLastEvaluatedKey)
		if err == nil {
			lastEvaluatedKeyString = string(marshalledKey)
		}
	}

	// Convert list to JSON
	messagelistJson, err := json.Marshal(messages)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// Put converted JSON in response
	fmt.Printf("Returning list of messages")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":   "*",
			"Access-Control-Allow-Methods":  "*",
			"Access-Control-Allow-Headers":  "*",
			"Access-Control-Expose-Headers": "Last-Evaluated-Key",
			"Content-Type":                  "application/json",
			"Last-Evaluated-Key":            lastEvaluatedKeyString,
		},
		Body: string(messagelistJson),
	}, nil
}

func main() {
	lambda.Start(getMessageHandler)
}
