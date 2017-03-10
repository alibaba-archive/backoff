package main

import (
	"net/http"
	"time"

	"github.com/teambition/backoff"
)

func main() {
	b := backoff.NewExponential()
	b.InitInterval = time.Second
	b.MaxInterval = 30 * time.Second
	b.SetMultiplier(1.5)
	b.SetFactor(0)
	b.Reset()

	var next time.Duration
	for {
		resp, err := http.DefaultClient.Get("http://www.baidu.com")
		if err == nil {
			// dosomething with resp
			_ = resp
			break
		}
		if next = b.Next(); next == backoff.Stop {
			break
		}
		time.Sleep(next)
	}
}
