package main

import (
	"fmt"
	"time"

	"github.com/hoodnoah/eve_market/monitor/dates"
	"github.com/hoodnoah/eve_market/monitor/downloader"
	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
)

const RequestsPerSecond = 3
const NumWorkers = 2

func fakeDateFn() time.Time {
	return time.Date(2003, 10, 5, 0, 0, 0, 0, time.UTC)
}

func main() {
	// setup logger
	logger := logger.NewConsoleLogger(10) // a healthy buffer since multiple services log
	logger.Start()

	// setup channels
	datesChannel := make(chan time.Time, NumWorkers)
	readerChannel := make(chan *parser.DatedReader, NumWorkers)
	marketDayChannel := make(chan *parser.MarketDay, NumWorkers)
	defer close(datesChannel)
	defer close(readerChannel)
	defer close(marketDayChannel)
	logger.Debug("channels initialized")

	// setup services
	dateIterator := dates.NewDateIterator(fakeDateFn, datesChannel)
	downloadManager := downloader.NewDownloadManager(logger, RequestsPerSecond, NumWorkers, readerChannel)
	parseManager := parser.NewMarketDataParser(logger, parser.DecompressFile, parser.ParseFile, readerChannel, marketDayChannel, NumWorkers)
	logger.Debug("services initialized")

	// start services
	dateIterator.Start()
	downloadManager.Start(datesChannel)
	parseManager.Start()
	logger.Debug("services started")

	for parsedResult := range marketDayChannel {
		msg := fmt.Sprintf("Parsed %d records for %s", len(parsedResult.Records), parsedResult.Date.Format(time.DateOnly))
		logger.Debug(msg)
	}
}
