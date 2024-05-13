package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}

	eventBridgeClient := eventbridge.New(session.New())

	eventInput := &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				Detail:       aws.String(string(event.Name)),
				DetailType:   aws.String("MessageCreation"),
				EventBusName: aws.String("KwetterEventBus"),
				Source:       aws.String("com.example.kwetter"),
			},
		},
	}

	message := fmt.Sprintf("Hello %s!", event.Name)

	_, err := eventBridgeClient.PutEvents(eventInput)
	if err != nil {
		return nil, fmt.Errorf("Error publishing event to EventBridge: %w", err)
	}
	return &message, nil
}

func main() {
	lambda.Start(HandleRequest)
}
