package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Post struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

func main() {
	url := "http://172.211.231.227:9000/users-api/Authentication/login"
	userID := "janou"
	content := "meh7"

	post := Post{UserName: userID, Password: content}
	postBody, _ := json.Marshal(post)
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{"application/json"}

	rate := vegeta.Rate{Freq: 300, Per: time.Second} // Adjust rate as needed
	duration := 1 * time.Minute                      // Adjust duration as needed
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
