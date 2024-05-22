package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

type Detail struct {
	Message string `json:"Message"`
	UserID  int    `json:"UserID"`
}

type Event struct {
	Source     string `json:"source"`
	Detail     Detail `json:"detail"`
	DetailType string `json:"detailType"`
}

func GenerateUUID() string {
	return uuid.New().String()
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
	currentTime := time.Now()
	timestamp := currentTime.Format(time.RFC3339)

	item := map[string]types.AttributeValue{
		"GUID":    &types.AttributeValueMemberS{Value: GenerateUUID()},
		"ts":      &types.AttributeValueMemberS{Value: timestamp},
		"Message": &types.AttributeValueMemberS{Value: event.Detail.Message},
		"UserID":  &types.AttributeValueMemberN{Value: strconv.Itoa(event.Detail.UserID)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("KwetterMessagesDynamoDBTable"),
		Item:      item,
	}

	_, err = client.PutItem(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
