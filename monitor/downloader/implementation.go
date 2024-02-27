package downloader

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
)

func newRateLimiter(requestsPerSecond uint) *RateLimiter {
	return &RateLimiter{
		ticker:      time.NewTicker(time.Second / time.Duration(requestsPerSecond)),
		tokenBucket: make(chan bool, requestsPerSecond),
	}
}

// exposes the rate limiter channel
func (rl *RateLimiter) StartRateLimiter() chan bool {
	// populate the bucket
	go func() {
		for range rl.ticker.C {
			select {
			case rl.tokenBucket <- true:
			default:
			}
		}
	}()

	return rl.tokenBucket
}

// constructor for a download manager
func NewDownloadManager(logger logger.ILogger, requestsPerSecond uint, numWorkers uint, resultsChannel chan *parser.DatedReader) *DownloadManager {
	if logger == nil {
		panic("Failed to create DownloadManager; logger must not be nil")
	}

	return &DownloadManager{
		rateLimiter:    *newRateLimiter(requestsPerSecond),
		numWorkers:     numWorkers,
		resultsChannel: resultsChannel,
		logger:         logger,
		mutex:          sync.Mutex{},
	}
}

func downloadFn(date time.Time) (*http.Response, error) {
	url := dateToURL(date)

	return http.Get(url)
}

func (dm *DownloadManager) handleError(err error, date time.Time) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	msg := fmt.Sprintf("Failed to download the file for %s: %v", date.Format(time.DateOnly), err)
	dm.logger.Error(msg)
}

// starts the downloading of dates
func (dm *DownloadManager) Start(datesChan chan time.Time) {
	rateChannel := dm.rateLimiter.StartRateLimiter()

	for range dm.numWorkers {
		go func() {
			for date := range datesChan {
				<-rateChannel
				res, err := downloadFn(date)

				if err != nil {
					dm.handleError(err, date)
				} else if res.StatusCode != http.StatusOK {
					errorMsg := fmt.Errorf("request received status code %d: %s", res.StatusCode, res.Status)
					dm.handleError(errorMsg, date)
				} else {
					dm.mutex.Lock()
					dm.resultsChannel <- &parser.DatedReader{
						Day:    date,
						Reader: res.Body,
					}
					dm.mutex.Unlock()
				}

			}
		}()
	}
}
