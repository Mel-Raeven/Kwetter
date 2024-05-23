package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var cognitoClient *cognitoidentityprovider.CognitoIdentityProvider

func init() {
	sess := session.Must(session.NewSession())
	cognitoClient = cognitoidentityprovider.New(sess)
}

func handler(event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// Extract the token from the AuthorizationToken field
	token := event.AuthorizationToken
	if token == "" {
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}

	// Extract the ID token from the cookie
	idToken := extractIDTokenFromCookie(token)
	if idToken == "" {
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}

	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: &idToken,
	}

	_, err := cognitoClient.GetUser(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to validate token:", err)
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}

	return generatePolicy("user", "Allow", event.MethodArn), nil
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{
		PrincipalID: principalID,
		PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		},
	}

	return authResponse
}

func extractIDTokenFromCookie(cookie string) string {
	cookies := strings.Split(cookie, "; ")
	for _, c := range cookies {
		if strings.HasPrefix(c, "id_token=") {
			return strings.TrimPrefix(c, "id_token=")
		}
	}
	return ""
}

func main() {
	lambda.Start(handler)
}
