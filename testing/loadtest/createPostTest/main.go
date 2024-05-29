package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Post struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

func main() {
	url := "https://6xoa9t5ole.execute-api.eu-central-1.amazonaws.com/Prod/postMessage"
	userID := "035408b2-a061-70bb-7f70-b1b4f150b359"
	content := "This is a test post"

	post := Post{UserID: userID, Content: content}
	postBody, _ := json.Marshal(post)
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{"application/json"}
	headers["Authorization"] = []string{"eyJraWQiOiJMK0ZMYklJWGxRMmViRno1eFlXV2dtNjBzQTNvRHFvcE93R1NJem9ZVkt3PSIsImFsZyI6IlJTMjU2In0.eyJhdF9oYXNoIjoiZDNOQUU0RTBfOGxUVVF3b0k4MGF0QSIsInN1YiI6IjAzNTQwOGIyLWEwNjEtNzBiYi03ZjcwLWIxYjRmMTUwYjM1OSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuZXUtY2VudHJhbC0xLmFtYXpvbmF3cy5jb21cL2V1LWNlbnRyYWwtMV9zRFVPNUhGd1giLCJjb2duaXRvOnVzZXJuYW1lIjoiMDM1NDA4YjItYTA2MS03MGJiLTdmNzAtYjFiNGYxNTBiMzU5Iiwib3JpZ2luX2p0aSI6IjRkMzY4ZTdhLTEwZmMtNGZiMy04YmMzLTRhZDg4ODhmNWExNCIsImF1ZCI6IjV2czRhbmpxdG1lNWY5MzAzN3Z2MzBmbXAwIiwidG9rZW5fdXNlIjoiaWQiLCJhdXRoX3RpbWUiOjE3MTY5NzIxODMsImV4cCI6MTcxNjk3NTc4MywiaWF0IjoxNzE2OTcyMTgzLCJqdGkiOiJiYjNmMGYxMy0yNDczLTQ0ZjYtYjNjMC1jOTllMTYwOTEwMmIiLCJlbWFpbCI6Im1lbC5yYWV2ZW5AZ21haWwuY29tIn0.HFM53OGfgIP66Ff11i5uZoNTeJ8GdsqW9_uQES9-Xtq-3GvKzD5luoxk5MKsTjEYkrITXj4oED1ToHqw6w2DB0GmiQl-PWtUE4-brULG1XGlNWHv26MkjBv43u5wdwo5YAAweZ5GVHMy1W_sSbrtyq6JFvOn_KJ5KSbrMANmenrbCh2u014WVB4XKz5gY7EdIGVnFV8fpnbuYEg3I2VWxKlX_lm7RGAFDqSLm0eetQb4gOUJaTVIrwtjo8cQyQygPdzHMgbcjNSJ1SxAqQPBRL7VtFPPnO0C_3ga1SduiHEX8EuHEd917PWtHqknME-tpyfMZq5h9ks38v8tj4buqw"}

	rate := vegeta.Rate{Freq: 10, Per: time.Second} // Adjust rate as needed
	duration := 2 * time.Minute                     // Adjust duration as needed
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
