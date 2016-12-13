package main

import (
	"net/http"
	"time"

	"github.com/xusss/backoff"
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
		_, err := http.DefaultClient.Get("http://www.baidu.com")
		if err == nil {
			break
		}
		if next = b.Next(); next == backoff.Stop {
			break
		}
		// dosomething and break
		time.Sleep(next)
	}
}
