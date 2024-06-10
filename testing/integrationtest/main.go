package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/eventbridge"
)

// TestEventBusMessageCreation tests the message creation event on a specific AWS Event Bus
func testEventBus() error {
	// Initialize AWS session using default credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"), // Replace "your-region" with your AWS region
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	// Initialize EventBridge client
	eventBridgeClient := eventbridge.New(sess)

	// Specify the name of your EventBridge bus
	eventBusName := "kwetter-ci-KwetterEventBus" // Replace "your-event-bus-name" with your EventBridge bus name

	// Generate a test message creation event
	event := &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				EventBusName: aws.String(eventBusName), // Specify the EventBridge bus name
				Source:       aws.String("com.example.kwetter"),
				Detail: aws.String(`{
					"detail-type": "MessageCreation",
					"source": "com.example.kwetter
					"detail": {
						"UserID": "user123",
						"Message": "Test message"
					}
				}`),
			},
		},
	}

	// Put the test event to the event bus
	_, err = eventBridgeClient.PutEvents(event)
	if err != nil {
		return fmt.Errorf("failed to put test event to event bus: %v", err)
	}

	// Sleep for a while to allow the event to be processed by your Lambda function
	time.Sleep(10 * time.Second) // Adjust sleep time as needed

	// Validate if the message has been created in DynamoDB
	dynamodbClient := dynamodb.New(sess)
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("kwetter-ci-KwetterMessagesDynamoDBTable"), // Replace with your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String("user123"),
			},
		},
	}

	_, err = dynamodbClient.GetItem(getItemInput)
	if err != nil {
		return fmt.Errorf("failed to get item from DynamoDB: %v", err)
	}

	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String("KwetterMessagesDynamoDBTable"), // Replace with your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				S: aws.String("user123"),
			},
		},
	}

	_, err = dynamodbClient.DeleteItem(deleteItemInput)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := testEventBus()
	if err != nil {
		fmt.Println(err)
	}
}
