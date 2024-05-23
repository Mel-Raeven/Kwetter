package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// TokenPayload represents the payload of a JWT token
type TokenPayload struct {
	Sub string `json:"sub"`
	// Add other fields as needed
}

// Handler is the Lambda function handler
func Handler(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// Log the received request for debugging
	fmt.Printf("Received request: %+v\n", request)

	// Retrieve the authorization token from the Authorization header
	authHeader := request.AuthorizationToken

	// If the token is empty, return Unauthorized
	if authHeader == "" {
		fmt.Println("Unauthorized: Missing Authorization header")
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized: Missing Authorization header")
	}

	// Split the Authorization header to get the token
	tokens := strings.Split(authHeader, " ")
	if len(tokens) != 2 || strings.ToLower(tokens[0]) != "bearer" {
		fmt.Printf("Unauthorized: Invalid Authorization header format. Received: %s\n", authHeader)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized: Invalid Authorization header format. Expected format: 'Bearer <token>'")
	}

	accessToken := tokens[1] // Extract the access token

	// Validate the access token
	accessTokenPayload, err := validateToken(accessToken)
	if err != nil {
		fmt.Println("Error validating access token:", err)
		return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized: Invalid token")
	}

	// Generate the policy document
	policyDocument := generatePolicy(accessTokenPayload)

	return policyDocument, nil
}

// validateToken validates the access token against Amazon Cognito
func validateToken(tokenString string) (*TokenPayload, error) {
	// If the token is empty, return an error
	if tokenString == "" {
		return nil, errors.New("Token is empty")
	}

	sess := session.Must(session.NewSession())
	cognitoIDP := cognitoidentityprovider.New(sess)

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(tokenString),
	}

	result, err := cognitoIDP.GetUser(input)
	if err != nil {
		fmt.Printf("Error from cognitoIDP.GetUser: %s\n", err.Error())
		return nil, err
	}

	// Extract the payload from the result
	payload := &TokenPayload{
		Sub: *result.Username,
		// Extract other attributes as needed
	}

	return payload, nil
}

// generatePolicy generates the API Gateway policy document
func generatePolicy(accessTokenPayload *TokenPayload) events.APIGatewayCustomAuthorizerResponse {
	// Define user information
	user := map[string]interface{}{
		"userId": accessTokenPayload.Sub,
		// Add other user attributes as needed
	}

	// Define the API Gateway policy document
	policyDocument := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: fmt.Sprintf("%v", user["userId"]),
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{"*"}, // Allow access to all API resources
				},
			},
		},
		Context: user, // Attach user information to the request context
	}

	return policyDocument
}

func main() {
	lambda.Start(Handler)
}
