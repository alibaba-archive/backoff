package main

import (
	"net/http"
	"time"

	"github.com/teambition/backoff"
)

func main() {
	b := backoff.NewExponential()
	b.MaxElapsed = 2 * time.Minute
	b.MaxInterval = 30 * time.Second
	b.SetFactor(0)
	b.Reset()

	var next time.Duration
	for {
		_, err := http.DefaultClient.Get("http://www.baidu.com")
		if next = b.Next(); next == backoff.Stop {
			break
		}
		// dosomething and break
		time.Sleep(next)
	}
}
