package main

import (
	"fmt"
	"log"
	"time"

	dbs "github.com/hoodnoah/eve_market/monitor/dbservice"
	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

func main() {
	config := dbs.ConfigVars{
		User:   "testuser",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "dbservice_test",
	}

	testRecord1 := mds.MarketHistoryCSVRecord{
		Date:       time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     18,
		Average:    10,
		Highest:    10,
		Lowest:     10,
		Volume:     340,
		OrderCount: 2,
	}

	testRecord2 := mds.MarketHistoryCSVRecord{
		Date:       time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     20,
		Average:    14,
		Highest:    14,
		Lowest:     14,
		Volume:     143,
		OrderCount: 1,
	}

	testRecords := []mds.MarketHistoryCSVRecord{testRecord1, testRecord2}

	dbService, err := dbs.NewDBService(&config)
	defer dbService.Close()

	if err != nil {
		log.Fatalf("Failed to open database: %s", err)
	}

	err = dbService.BulkAddRecord(testRecords)

	if err != nil {
		log.Fatalf("Failed to insert record: %s", err)
	} else {
		fmt.Println("Successfully inserted record.")
	}
}
