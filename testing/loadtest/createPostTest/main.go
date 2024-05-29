package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Detail struct {
	Message string `json:"Message"`
	UserID  string `json:"UserID"`
}
type Post struct {
	Detail Detail `json:"detail"`
}

func main() {
	url := "https://0a43x0s4q4.execute-api.eu-central-1.amazonaws.com/Prod/postMessage"
	userID := "136428f2-8051-7099-d7a2-465acc1593bc"
	message := "This is a test post"

	post := Post{Detail: Detail{Message: message, UserID: userID}}
	postBody, _ := json.Marshal(post)
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{"application/json"}
	headers["Authorization"] = []string{"eyJraWQiOiJMUlZIWjJ6SU9jSkwrSDJmcHhmREg0UHJpbXVJQWdHVUgyM2hqMWpOVVRFPSIsImFsZyI6IlJTMjU2In0.eyJhdF9oYXNoIjoiT1VvT1BRaS0ySDZ4Sy1wSE5hbGhvdyIsInN1YiI6IjEzNjQyOGYyLTgwNTEtNzA5OS1kN2EyLTQ2NWFjYzE1OTNiYyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuZXUtY2VudHJhbC0xLmFtYXpvbmF3cy5jb21cL2V1LWNlbnRyYWwtMV9qeERPSTBPOEQiLCJjb2duaXRvOnVzZXJuYW1lIjoiMTM2NDI4ZjItODA1MS03MDk5LWQ3YTItNDY1YWNjMTU5M2JjIiwib3JpZ2luX2p0aSI6ImVlNzQyMTM5LWQ4ZDgtNDE5ZC1hMzYxLWVlODk0OGE0OTc1OSIsImF1ZCI6IjF2bmxkNXZsN3RvY21lY3Zia2prc21zZjJyIiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE3MTY5ODAzMjYsImV4cCI6MTcxNjk4MzkyNiwiaWF0IjoxNzE2OTgwMzI2LCJqdGkiOiJkZDliN2IwZi02OGY4LTRlNTctOWE2Yy1lMGIzYzRlMDVlNDkiLCJlbWFpbCI6Im1lbC5yYWV2ZW5AZ21haWwuY29tIn0.kR5QWrQnQy5w9aN6DY_LD3g2VupEfdGHU_zdOIHnZbWv21tJiV_RX4xHLYdRqroKN8rNo6bK4V7AN4vaZ8-w6idEfRMfLLFKBNHTf9TLxSHxU0uonr8RxB_I-nAo4BzV5R4HMv3zDhPf7cRbVWTQm10KX0vU9gK-SPgHuj6petGa8aCiysYoF0C9iFNcuI8ktKb0pfu0XJhNUPJVhyFrx1xIJG6jGe5hXetz909y2qan38g60z4JdZqarsvFBZN3lr3p2jMhBdRCJoX_yXwt9MiGi2pbENPEoP0DEzp6uctk-hL36plEPwGgTkokcHn81-Dp133Ig133mHT-zkRHzA"}

	rate := vegeta.Rate{Freq: 100, Per: time.Second} // Adjust rate as needed
	duration := 2 * time.Minute                      // Adjust duration as needed
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    url,
		Body:   postBody,
		Header: headers,
	})

	fmt.Println("Starting the attack")
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	ticker := time.NewTicker(1 * time.Second) // Adjust the interval as needed
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				printMetrics(metrics)
			}
		}
	}()

	for res := range attacker.Attack(targeter, rate, duration, "CreatePostTest") {
		metrics.Add(res)
	}
	metrics.Close()
	done <- true
	ticker.Stop()

	printMetrics(metrics)
}

func printMetrics(metrics vegeta.Metrics) {
	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
}
