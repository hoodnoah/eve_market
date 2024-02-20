package integration_tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/hoodnoah/eve_market/monitor/dbservice"
	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

func TestInsertOne(t *testing.T) {
	// configure connection to the test db
	config := mysql.Config{
		User:                 "testuser",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "dbservice_test",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// initialize the service, which will create the tables
	dbservice, err := db.NewMySqlDBService(&config)
	if err != nil {
		t.Fatalf("Failed to create MySqlDBService: %v", err)
	}
	defer dbservice.Close()

	testRecord := mds.MarketHistoryCSVRecord{
		Date:       time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     39,
		Average:    1060.2,
		Highest:    1100,
		Lowest:     1024,
		Volume:     41,
		OrderCount: 8,
	}

	err = dbservice.InsertOne(&testRecord)
	if err != nil {
		t.Fatalf("Failed to insert test record: %v", err)
	}

	// open a separate connection to test values
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		t.Fatalf("Failed to connect independently to the database: %v", err)
	}
	defer db.Close()

	var resultRecord mds.MarketHistoryCSVRecord

	err = db.QueryRow("SELECT * FROM market_data WHERE region_id = 10000001 LIMIT 1").Scan(
		&resultRecord.Date,
		&resultRecord.RegionID,
		&resultRecord.TypeID,
		&resultRecord.Average,
		&resultRecord.Highest,
		&resultRecord.Lowest,
		&resultRecord.Volume,
		&resultRecord.OrderCount,
	)
	if err != nil {
		t.Fatalf("Failed to query record: %v", err)
	}

	if !testRecord.Equals(&resultRecord) {
		t.Fatalf("Expected records to be equal: \n\nExpected:\n\n%v\n\nReceived:\n\n%v", testRecord, resultRecord)
	}
}
