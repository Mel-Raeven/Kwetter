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
	headers["Authorization"] = []string{"eyJraWQiOiJMUlZIWjJ6SU9jSkwrSDJmcHhmREg0UHJpbXVJQWdHVUgyM2hqMWpOVVRFPSIsImFsZyI6IlJTMjU2In0.eyJhdF9oYXNoIjoieHRZUzZ0TkoyVkRmVkNZNGxLN1M4ZyIsInN1YiI6IjEzNjQyOGYyLTgwNTEtNzA5OS1kN2EyLTQ2NWFjYzE1OTNiYyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuZXUtY2VudHJhbC0xLmFtYXpvbmF3cy5jb21cL2V1LWNlbnRyYWwtMV9qeERPSTBPOEQiLCJjb2duaXRvOnVzZXJuYW1lIjoiMTM2NDI4ZjItODA1MS03MDk5LWQ3YTItNDY1YWNjMTU5M2JjIiwib3JpZ2luX2p0aSI6IjU2NWVmYmQyLTVlZDItNDgzNS05YzlkLTE1NWZjYTVlMTgzNiIsImF1ZCI6IjF2bmxkNXZsN3RvY21lY3Zia2prc21zZjJyIiwiZXZlbnRfaWQiOiIzOTViNjYxNC03OGZiLTQ3ZWMtYmNiYi1jYTE1MWY3MGVmYzkiLCJ0b2tlbl91c2UiOiJpZCIsImF1dGhfdGltZSI6MTcxNjk5OTc2MiwiZXhwIjoxNzE3MDAzMzYyLCJpYXQiOjE3MTY5OTk3NjIsImp0aSI6ImI4MmUxZGZlLTY0YTgtNGUwZC1hMzE2LTdhZDE0MzA2NzRkZCIsImVtYWlsIjoibWVsLnJhZXZlbkBnbWFpbC5jb20ifQ.Qp7jdSESQS4xcij_6D91roC2q1PPLYCL6HkJvJIj-3penLjQGmoJ7T18W2uSbz9JjolhCL2RVYdcOuY-wskTZZnm6_hYPa0Lk8J3loFkOT92ChBKp4QyXixZfqrRn4t3b1a-S_8tJO_rZFlYqDZ8eImzj-JExD8g9dzDGm5Xrx81jLNGMrjCc-xhcqIojfDpKIZfLzUftWhZo8C56Xys_2QXMuQNnJ_oAYTw9x082lkKaGy7H-QMFOK6HFPPratVFbo5bp4mhqzqpzkS-1ocPlTm1s_ujs-wlUDJcL8szpIi05Xj7iG1ZNSiMPNQA1NmvX28tAgAfgXiDcMTTSh7FA"}

	rate := vegeta.Rate{Freq: 100, Per: time.Second} // Adjust rate as needed
	duration := 10 * time.Minute                     // Adjust duration as needed
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
