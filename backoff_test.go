package backoff

import (
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

func TestMaxRetry(t *testing.T) {
	maxRetry := 10
	b := NewExponentialWithRetry(maxRetry)
	b.Reset()

	var (
		next  time.Duration
		retry int
	)
	for {
		if next = b.Next(); next == Stop {
			break
		}
		retry += 1
	}
	if retry != 10 {
		t.Fatal("retry: ", retry)
	}
}

func TestMaxElapsed(t *testing.T) {
	maxElapsed := 10 * time.Second
	b := NewExponentialWithElapsed(maxElapsed)
	b.SetMultiplier(1)
	b.Reset()

	var (
		next time.Duration
	)
	start := time.Now()
	for {
		if next = b.Next(); next == Stop {
			break
		}
	}
	end := time.Now()
	// with tiny gap
	t.Logf("maxElapsed: %s\n", maxElapsed)
	t.Logf("elapsed: %s\n", end.Sub(start))
}
