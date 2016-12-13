package main

import (
	"fmt"
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
	for {
		next := b.Next()
		if next == backoff.Stop {
			break
		}
		fmt.Println(time.Now())
		_, err := http.DefaultClient.Get("http://www.baidu.com")
		if err != nil {
			continue
		}
		// dosomething and break
		time.Sleep(next)
	}
}
