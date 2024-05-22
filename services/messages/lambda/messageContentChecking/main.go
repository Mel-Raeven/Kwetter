package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
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

func HandleRequest(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("received nil event")
	}

	eventBridgeClient := eventbridge.New(session.New())

	eventData := map[string]interface{}{
		"Message": event.Detail.Message,
		"UserID":  event.Detail.UserID,
	}

	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	eventInput := &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				Source:       aws.String("com.example.kwetter"),
				Detail:       aws.String(string(jsonData)),
				DetailType:   aws.String("MessageCreation"),
				EventBusName: aws.String("KwetterEventBus"),
			},
		},
	}

	_, putEventErr := eventBridgeClient.PutEvents(eventInput)
	if putEventErr != nil {
		return putEventErr
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
