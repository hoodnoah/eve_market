package dates

import (
	"sync"
	"time"
)

type NowFn func() time.Time

type DateIterator struct {
	currentDate time.Time
	mutex       sync.Mutex
	nowFn       NowFn
}
