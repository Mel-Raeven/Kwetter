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
	ts      string `json:"ts"`
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

	// create dynamodb client
	client := dynamodb.NewFromConfig(cfg)

	// create scaninput that defines the conditions and tablename aswell as the expected returning fields
	input := &dynamodb.ScanInput{
		TableName:        aws.String("KwetterMessagesDynamoDBTable"),
		FilterExpression: aws.String("UserID = :userid"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userid": &types.AttributeValueMemberS{Value: id},
		},
		ProjectionExpression: aws.String("GUID, Message, ts, UserID"),
	}

	// perform the scan
	result, err := client.Scan(ctx, input)
	if err != nil {
		fmt.Printf("Couldn't scan info for user %v. Here's why: %v\n", id, err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	if len(result.Items) == 0 {
		fmt.Printf("No bookings found for user %v\n", id)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "*",
				"Access-Control-Allow-Headers": "*",
				"Content-Type":                 "application/json",
			},
			Body: "",
		}, nil
	}

	// create list of bookings from results
	var messages []Message
	err = attributevalue.UnmarshalListOfMaps(result.Items, &messages)
	if err != nil {
		fmt.Printf("Couldn't unmarshal scan response. Here's why: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// convert list to json
	messagelistJson, err := json.Marshal(messages)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, err
	}

	// put converted json in response
	fmt.Printf("Returning list of bookings")
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
			"Access-Control-Allow-Headers": "*",
			"Content-Type":                 "application/json",
		},
		Body: string(messagelistJson),
	}, nil
}

func main() {
	lambda.Start(getMessageHandler)
}
