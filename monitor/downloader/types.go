package downloader

import (
	"sync"
	"time"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
)

type RateLimiter struct {
	ticker      *time.Ticker
	tokenBucket chan bool
}

type DownloadManager struct {
	logger         logger.ILogger
	rateLimiter    RateLimiter
	datesChannel   chan time.Time
	resultsChannel chan *parser.DatedReader
	numWorkers     uint
	mutex          sync.Mutex
}
