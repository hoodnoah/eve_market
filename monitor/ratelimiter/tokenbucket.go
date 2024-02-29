package ratelimiter

import (
	"sync"
	"time"
)

func NewTokenBucketRateLimiter(requestsPerSecond int) IRateLimiter {
	return &TokenBucketRateLimiter{
		ticker:      time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
		tokenBucket: make(chan bool, requestsPerSecond),
		mutex:       sync.Mutex{},
	}
}

func (tb *TokenBucketRateLimiter) Start() {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	go func() {
		for range tb.ticker.C {
			select {
			case tb.tokenBucket <- true: // populate when ticker is ready
			default: // do nothing otherwise
			}
		}
	}()
}

func (tb *TokenBucketRateLimiter) GetChannel() chan bool {
	return tb.tokenBucket
}
