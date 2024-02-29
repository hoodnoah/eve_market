package downloader

import (
	"sync"
	"time"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
	"github.com/hoodnoah/eve_market/monitor/ratelimiter"
)

type DownloadManager struct {
	logger         logger.ILogger
	rateLimiter    ratelimiter.IRateLimiter
	datesChannel   chan time.Time
	resultsChannel chan *parser.DatedReader
	numWorkers     uint
	excludeDates   map[time.Time]bool
	mutex          sync.Mutex
}
