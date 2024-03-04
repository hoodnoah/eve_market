package downloader

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
	"github.com/hoodnoah/eve_market/monitor/ratelimiter"
)

// constructor for a download manager
func NewDownloadManager(logger logger.ILogger, requestsPerSecond uint, numWorkers uint, inputChannel chan time.Time, resultsChannel chan *parser.DatedReader) *DownloadManager {
	if logger == nil {
		panic("Failed to create DownloadManager; logger must not be nil")
	}

	rl := ratelimiter.NewTokenBucketRateLimiter(int(requestsPerSecond))
	rl.Start()

	return &DownloadManager{
		rateLimiter:    rl,
		numWorkers:     numWorkers,
		datesChannel:   inputChannel,
		resultsChannel: resultsChannel,
		logger:         logger,
		excludeDates:   make(map[time.Time]bool),
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

func (dm *DownloadManager) isExcludableDate(date time.Time) bool {
	return dm.excludeDates[date]
}

// starts the downloading of dates
func (dm *DownloadManager) Start() {
	rateChannel := dm.rateLimiter.GetChannel()

	for range dm.numWorkers {
		go func() {
			for date := range dm.datesChannel {
				// make sure the date isn't already retrieved
				if dm.isExcludableDate(date) {
					continue
				}

				// block on rate limited
				<-rateChannel

				// execute download
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

func (dm *DownloadManager) Exclude(dates []time.Time) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for _, date := range dates {
		dm.excludeDates[date] = true
	}
}
