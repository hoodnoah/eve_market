package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hoodnoah/eve_market/monitor/dateservice"
)

func GetRateLimiter(requestsPerSecond int) chan bool {
	interval := time.Second / time.Duration(requestsPerSecond)
	ticker := time.NewTicker(interval)

	tokenBucket := make(chan bool, requestsPerSecond)

	// start replenishing the tokenbucket
	go func() {
		for range ticker.C {
			select {
			case tokenBucket <- true:
				// token added
			default:
				// bucket full
			}
		}
	}()

	return tokenBucket
}

func mockDownloader(thing int) {
	fmt.Printf("Downloading item %d...\n", thing)
	time.Sleep(2 * time.Second)
}

func main() {
	// setup services
	dateSvc := dateservice.NewEveRefsDateService(time.Now)

	// get dates
	dates := dateSvc.EnumerateDates()

	fmt.Printf("Number of dates between 2003-10-01 and now: %d\n\n", len(dates))

	thingsToDownload := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	downloaderChannel := make(chan int, len(thingsToDownload))
	for _, dl := range thingsToDownload {
		downloaderChannel <- dl
	}

	tokenBucket := GetRateLimiter(2)
	var wg sync.WaitGroup

	startTime := time.Now()
	for i := 1; i <= 2; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for item := range downloaderChannel {
				<-tokenBucket
				mockDownloader(item)
			}
		}()
	}

	close(downloaderChannel)

	wg.Wait()
	endTime := time.Now()

	timeElapsed := endTime.Sub(startTime)
	fmt.Println(timeElapsed)
}
