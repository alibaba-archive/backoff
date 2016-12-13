package backoff

import (
	"math/rand"
	"time"
)

/*
simple backoff algorithm implement with golang

example:
	backoff.Reset()

	var next time.Duration
	for {
		// do something
		if err == nil {
			break
		}
		if next = backoff.Next(); next == backoff.Stop {
			break
		}
		time.Sleep(next)
	}
*/

// BackOff algorithm for calculate the period of each retry
type BackOff interface {
	Next() time.Duration
	// call Reset() before use
	Reset()
}

const Stop time.Duration = -1

const (
	DefaultFactor       = 0.5
	DefaultInitInterval = 500 * time.Millisecond
	DefaultMultiplier   = 1.5
	DefaultMaxInterval  = 60 * time.Second
)

// Exponential with the exponential growth period for each retry
type Exponential struct {
	// need currentInterval * multiplier after calculate the NextInterval
	// factor must in [0, 1]
	factor, multiplier float64
	// if now - startTime > MaxElapsed then stop
	// if the growth period reach the MaxInterval/Multiplier then return the MaxInterval
	InitInterval, currentInterval, MaxInterval, maxElapsed time.Duration
	reachMaxInterval                                       bool
	startTime                                              time.Time
	// the number of retries
	maxRetry, currentRetry int
}

func NewExponentialWithElapsed(elapsed time.Duration) *Exponential {
	b := NewExponential()
	b.maxElapsed = elapsed
	return b
}

func NewExponentialWithRetry(retry int) *Exponential {
	b := NewExponential()
	b.maxRetry = retry
	return b
}

func NewExponential() *Exponential {
	b := &Exponential{
		factor:       DefaultFactor,
		multiplier:   DefaultMultiplier,
		InitInterval: DefaultInitInterval,
		MaxInterval:  DefaultMaxInterval,
	}
	return b
}

// nextInterval = currentInterval * (random between [1 + Factor, 1 - Factor])
func (b *Exponential) nextInterval() time.Duration {
	rnd := rand.Float64()
	max := float64(b.currentInterval) * (1 + b.factor)
	min := float64(b.currentInterval) * (1 - b.factor)

	return time.Duration(min + rnd*(max-min+1))
}

func (b *Exponential) incrCurrent() {
	if b.reachMaxInterval {
		return
	}
	maybe := float64(b.currentInterval) * b.multiplier
	if maybe >= float64(b.MaxInterval) {
		b.currentInterval = b.MaxInterval
		b.reachMaxInterval = true
	} else {
		b.currentInterval = time.Duration(maybe)
	}
}

func (b *Exponential) Next() time.Duration {
	if b.maxElapsed != 0 && time.Now().Sub(b.startTime) > b.maxElapsed {
		return Stop
	}
	if b.maxRetry != 0 && b.currentRetry >= b.maxRetry {
		return Stop
	}
	b.currentRetry += 1
	defer b.incrCurrent()
	return b.nextInterval()
}

func (b *Exponential) Reset() {
	b.currentInterval = b.InitInterval
	b.startTime = time.Now()
}

func (b *Exponential) SetFactor(factor float64) {
	switch {
	case factor < 0:
		b.factor = 0
	case factor > 1:
		b.factor = 1
	default:
		b.factor = factor
	}
}

func (b *Exponential) Factor() float64 {
	return b.factor
}

func (b *Exponential) SetMultiplier(multiplier float64) {
	switch {
	case multiplier < 1:
	default:
		b.multiplier = multiplier
	}
}

func (b *Exponential) Multiplier() float64 {
	return b.multiplier
}
