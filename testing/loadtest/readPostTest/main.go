package main

import (
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	url := "http://your-aws-app-url/api/posts"
	rate := vegeta.Rate{Freq: 1000, Per: time.Second} // Adjust rate as needed
	duration := 5 * time.Minute                       // Adjust duration as needed
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    url,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "FetchPostsTest") {
		metrics.Add(res)
	}
	metrics.Close()

	reporter := vegeta.NewTextReporter(&metrics)
	reporter.Report(os.Stdout)
}
