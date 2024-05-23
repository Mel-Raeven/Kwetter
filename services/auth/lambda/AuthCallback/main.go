package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ClaimsOverrideDetails struct {
	ClaimsToAddOrOverride map[string]string `json:"claimsToAddOrOverride,omitempty"`
}

type CognitoEvent struct {
	Version               string                 `json:"version"`
	TriggerSource         string                 `json:"triggerSource"`
	Region                string                 `json:"region"`
	UserPoolID            string                 `json:"userPoolId"`
	CallerContext         map[string]interface{} `json:"callerContext"`
	UserName              string                 `json:"userName"`
	Request               map[string]interface{} `json:"request"`
	Response              map[string]interface{} `json:"response"`
	ClaimsOverrideDetails `json:"claimsOverrideDetails,omitempty"`
}

func handler(event CognitoEvent) (events.APIGatewayProxyResponse, error) {
	// Log the incoming event for debugging purposes
	eventJSON, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	fmt.Printf("Received event: %s\n", string(eventJSON))

	// Generate the HTTP-only cookie
	idToken, ok := event.Request["userAttributes"].(map[string]interface{})["sub"].(string)
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message": "Invalid user attributes"}`,
		}, nil
	}
	cookie := fmt.Sprintf("id_token=%s; HttpOnly; Secure; SameSite=None; Path=/; Max-Age=3600", idToken)

	// Create the response
	headers := map[string]string{
		"Set-Cookie":   cookie,
		"Content-Type": "application/json",
	}
	body, err := json.Marshal(map[string]string{
		"message": "Token set in HttpOnly cookie",
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       `{"message": "Failed to create response body"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
