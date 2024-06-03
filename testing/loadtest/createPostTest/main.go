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
	userID := "93f45812-4051-706c-61ae-c7078b86e4bb"
	message := "This is a test post"

	post := Post{Detail: Detail{Message: message, UserID: userID}}
	postBody, _ := json.Marshal(post)
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{"application/json"}
	headers["Authorization"] = []string{"eyJraWQiOiJMUlZIWjJ6SU9jSkwrSDJmcHhmREg0UHJpbXVJQWdHVUgyM2hqMWpOVVRFPSIsImFsZyI6IlJTMjU2In0.eyJhdF9oYXNoIjoiSmh6bjVNOFBqN2I3LXl6a2ltbVpfZyIsInN1YiI6IjEzNjQyOGYyLTgwNTEtNzA5OS1kN2EyLTQ2NWFjYzE1OTNiYyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuZXUtY2VudHJhbC0xLmFtYXpvbmF3cy5jb21cL2V1LWNlbnRyYWwtMV9qeERPSTBPOEQiLCJjb2duaXRvOnVzZXJuYW1lIjoiMTM2NDI4ZjItODA1MS03MDk5LWQ3YTItNDY1YWNjMTU5M2JjIiwib3JpZ2luX2p0aSI6IjJmZDg2NzQ5LThmZTktNGVmZC04NDNlLWJiZTQ1ZDNjNDdmNyIsImF1ZCI6IjF2bmxkNXZsN3RvY21lY3Zia2prc21zZjJyIiwiZXZlbnRfaWQiOiJhZDBlZjgyOC03OWVhLTQ3ZTItOTcyMi0yNzg5ODc1NjRkZTYiLCJ0b2tlbl91c2UiOiJpZCIsImF1dGhfdGltZSI6MTcxNzQxNzIzNCwiZXhwIjoxNzE3NDIwODM0LCJpYXQiOjE3MTc0MTcyMzQsImp0aSI6IjRjMzY4MDcxLTE1NmQtNGRkNC1hNzQwLTM0OWQ1NDY2MDVlMiIsImVtYWlsIjoibWVsLnJhZXZlbkBnbWFpbC5jb20ifQ.GdRm3Vseoha873CPypbLyy-xTa9BztvoFfQwU3lBGiOCMiRKKDr4gEHKWGyZskNiAh3n4uNrTv0CfnoawB2NnUApvr80ZLG5sAwA5cOA1lFYkcNTzB0Zg_wMFmZ3ZCWRtDtkt7-UJ2sz4mKS7tOCRarcXhOeU2t6pYZR3w_KMx95_An-L5CfkY8z_qf7wZni5g6VlluJ2vjX1s5xT5Mqgpr8in7oO7yd0D1DsCqruk0P8KBdIR28FxnsPyGhKEqoY99DCAqSUetg78XoFP5uWnQcJ_B1ezZTRIl0orAUeiLOhBjPqQN6tHLp1iRX7WlBMG0WlftMAaCNlbjrR9jeTw"}

	rate := vegeta.Rate{Freq: 50, Per: time.Second} // Adjust rate as needed
	duration := 1 * time.Minute                     // Adjust duration as needed
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
		// fmt.Printf("Request timestamp: %s\n", time.Now().Format(time.RFC3339))
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
