package main

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hoodnoah/eve_market/monitor/dates"
	"github.com/hoodnoah/eve_market/monitor/dbservice"
	"github.com/hoodnoah/eve_market/monitor/downloader"
	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
)

const RequestsPerSecond = 3
const NumWorkers = 2

func fakeDateFn() time.Time {
	return time.Date(2003, 10, 5, 0, 0, 0, 0, time.UTC)
}

var MySqlConfig mysql.Config = mysql.Config{
	User:                 "testuser",
	Passwd:               "password",
	Net:                  "tcp",
	Addr:                 "localhost:3306",
	DBName:               "dbservice_test",
	AllowNativePasswords: true,
	ParseTime:            true,
}

func main() {
	// setup logger
	logger := logger.NewConsoleLogger(10) // a healthy buffer since multiple services log
	logger.Start()

	// setup channels
	datesChannel := make(chan time.Time, NumWorkers)
	readerChannel := make(chan *parser.DatedReader, NumWorkers)
	marketDayChannel := make(chan *parser.MarketDay, NumWorkers)
	dbResultsChannel := make(chan time.Time, NumWorkers)
	defer close(datesChannel)
	defer close(readerChannel)
	defer close(marketDayChannel)
	defer close(dbResultsChannel)
	logger.Debug("channels initialized")

	// setup services
	dateIterator := dates.NewDateIterator(fakeDateFn, datesChannel)
	downloadManager := downloader.NewDownloadManager(logger, RequestsPerSecond, NumWorkers, readerChannel)
	parseManager := parser.NewMarketDataParser(logger, parser.DecompressFile, parser.ParseFile, readerChannel, marketDayChannel, NumWorkers)
	dbManager, err := dbservice.NewDBManager(&MySqlConfig, logger, marketDayChannel, dbResultsChannel, NumWorkers)
	if err != nil {
		msg := fmt.Sprintf("Failed to initalize a new dbmanager: %s", err)
		panic(msg)
	}
	defer dbManager.Close()
	logger.Debug("services initialized")

	// start services
	dateIterator.Start()
	downloadManager.Start(datesChannel)
	parseManager.Start()
	dbManager.Start()
	logger.Debug("services started")

	for dbResult := range dbResultsChannel {
		msg := fmt.Sprintf("Successfully inserted record for %s", dbResult.Format(time.DateOnly))
		logger.Debug(msg)
	}
}
