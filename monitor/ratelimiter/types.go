package ratelimiter

import (
	"sync"
	"time"
)

type IRateLimiter interface {
	Start()
	GetChannel() chan bool
}

type TokenBucketRateLimiter struct {
	ticker      *time.Ticker
	tokenBucket chan bool
	mutex       sync.Mutex
}
