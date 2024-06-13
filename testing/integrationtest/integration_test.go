package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/stretchr/testify/require"
)

func getPhysicalID(svc *cloudformation.CloudFormation, stackName, logicalID string) (string, error) {
	input := &cloudformation.DescribeStackResourcesInput{
		StackName:         aws.String(stackName),
		LogicalResourceId: aws.String(logicalID),
	}

	result, err := svc.DescribeStackResources(input)
	if err != nil {
		return "", fmt.Errorf("failed to describe stack resources: %v", err)
	}

	if len(result.StackResources) == 0 {
		return "", fmt.Errorf("no resources found with the provided logical ID")
	}

	return *result.StackResources[0].PhysicalResourceId, nil
}

func testMessageIntegration(t *testing.T) {
	// Initialize AWS session using default credentials
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"), // Replace with your AWS region
	})
	require.NoError(t, err, "failed to create session")

	// stackname to get the dynamodb table name
	stackName := "kwetter-ci"                   // Replace with your actual stack name
	logicalID := "KwetterMessagesDynamoDBTable" // Replace with the logical ID of your DynamoDB table
	cfClient := cloudformation.New(sess)

	// Initialize EventBridge client
	eventBridgeClient := eventbridge.New(sess)

	// Specify the name of your EventBridge bus
	eventBusName := "kwetter-ci-KwetterEventBus" // Replace with your EventBridge bus name

	// Generate a test message creation event
	event := &eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				EventBusName: aws.String(eventBusName), // Specify the EventBridge bus name
				Source:       aws.String("com.example.kwetter"),
				DetailType:   aws.String("MessageCreation"),
				Detail: aws.String(`{
					"UserID": "user123",
					"Message": "Test message"
				}`),
			},
		},
	}

	// Put the test event to the event bus
	_, err = eventBridgeClient.PutEvents(event)
	require.NoError(t, err, "failed to put test event to event bus")

	// Sleep for a while to allow the event to be processed by your Lambda function
	time.Sleep(10 * time.Second) // Adjust sleep time as needed

	// Validate if the message has been created in DynamoDB
	dynamodbClient := dynamodb.New(sess)

	tableName, err := getPhysicalID(cfClient, stackName, logicalID)
	require.NoError(t, err, "failed to get physical ID")

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userID": {
				S: aws.String("user123"),
			},
		},
	}

	result, err := dynamodbClient.Query(queryInput)
	require.NoError(t, err, "failed to query item from DynamoDB")
	require.NotEmpty(t, result.Items, "no item found with the provided UserID")

	// Optional: Add more assertions to validate the retrieved item
}

func TestMessageIntegration(t *testing.T) {
	testMessageIntegration(t)
}
