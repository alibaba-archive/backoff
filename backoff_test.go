package backoff

import (
	"fmt"
	"testing"
	"time"
)

func TestBoundOfFactor(t *testing.T) {
	b := NewExponential()
	b.SetFactor(-1)
	if b.Factor() != 0 {
		t.Fatal("b.factor: ", b.Factor())
	}
	b.SetFactor(3)
	if b.Factor() != 1 {
		t.Fatal("b.factor: ", b.Factor())
	}
}

func TestNext(t *testing.T) {
	b := NewExponential()
	b.MaxRetry = 10
	b.Reset()

	var next time.Duration
	for {
		if next = b.Next(); next == Stop {
			break
		}
		fmt.Printf("period:  %s\n", next)
		time.Sleep(time.Second)
	}
}
