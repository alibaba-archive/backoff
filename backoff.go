package backoff

import (
	"fmt"
	"math/rand"
	"time"
)

/*
simple backoff algorithm implement with golang

example:
	backoff.Reset()
	for {
		next := backoff.Next()
		if next == backoff.Stop {
			break
		}
		// do something
		time.Sleep(next)
	}
*/

// BackOff algorithm for calculate the period of each retry
type BackOff interface {
	Next() time.Duration
	Reset()
}

const Stop time.Duration = -1

const (
	DefaultFactor       = 0.5
	DefaultInitInterval = 500 * time.Millisecond
	DefaultMultiplier   = 1.5
	DefaultMaxInterval  = 60 * time.Second
	DefaultMaxElapsed   = 15 * time.Minute
	DefaultMaxRetry     = 30
)

// Exponential with the exponential growth period for each retry
type Exponential struct {
	// need currentInterval * multiplier after calculate the NextInterval
	// factor must in [0, 1]
	factor, multiplier float64
	// if now - startTime > MaxElapsed then stop
	// if the growth period reach the MaxInterval/Multiplier then return the MaxInterval
	InitInterval, currentInterval, MaxInterval, MaxElapsed time.Duration
	reachMaxInterval                                       bool
	startTime                                              time.Time
	// the number of retries
	MaxRetry, currentRetry int
}

func NewExponential() *Exponential {
	b := &Exponential{
		factor:       DefaultFactor,
		multiplier:   DefaultMultiplier,
		InitInterval: DefaultInitInterval,
		MaxInterval:  DefaultMaxInterval,
		MaxElapsed:   DefaultMaxElapsed,
		MaxRetry:     DefaultMaxRetry,
	}
	b.Reset()
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
	if b.MaxElapsed != 0 && time.Now().Sub(b.startTime) > b.MaxElapsed {
		return Stop
	}
	fmt.Println(b.currentRetry)
	if b.MaxRetry != 0 && b.currentRetry >= b.MaxRetry {
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
	case multiplier <= 1:
	default:
		b.multiplier = multiplier
	}
}

func (b *Exponential) Multiplier() float64 {
	return b.multiplier
}
