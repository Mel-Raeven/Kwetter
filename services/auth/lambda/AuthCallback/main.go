// main.go
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	code := event.QueryStringParameters["code"]
	if code == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Authorization code not found",
		}, nil
	}

	tokenResponse, err := fetchTokens(code)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error fetching tokens",
		}, nil
	}

	cookies := make([]string, 3)
	cookies[0] = createCookie("access_token", tokenResponse.AccessToken)
	cookies[1] = createCookie("id_token", tokenResponse.IdToken)
	cookies[2] = createCookie("refresh_token", tokenResponse.RefreshToken)

	responseHeaders := map[string]string{
		"Set-Cookie":                  strings.Join(cookies, "; "),
		"Location":                    "http://localhost:5173/auth/callback",
		"Access-Control-Allow-Origin": "*",
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 302,
		Headers:    responseHeaders,
	}, nil
}

func fetchTokens(code string) (*TokenResponse, error) {
	clientID := "your-client-id"
	clientSecret := "your-client-secret"
	redirectURI := "https://your-service.com/auth/callback"
	tokenURL := "https://your-cognito-domain/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func createCookie(name, value string) string {
	return name + "=" + value + "; HttpOnly; Secure; Path=/"
}
