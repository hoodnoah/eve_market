package integration_tests

import (
	"bytes"
	"os"
	"testing"
	"time"

	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

func TestInsertMany(t *testing.T) {
	t.Run("inserts a single sample record correctly", func(t *testing.T) {
		setup, err := Setup()
		if err != nil {
			t.Fatalf("Failed to setup test: %v", err)
		}
		defer setup.TearDown()

		testRecord := []mds.MarketHistoryCSVRecord{
			{
				Date:       time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC),
				RegionID:   1,
				TypeID:     1,
				Average:    1,
				Highest:    2,
				Lowest:     1,
				Volume:     1,
				OrderCount: 1,
			},
		}

		err = (*setup.DBService).InsertMany(testRecord)
		if err != nil {
			t.Fatalf("Failed to InsertMany a single record: %v", err)
		}

		// get record
		row := setup.Connection.QueryRow(`
			SELECT * FROM market_data;
			`)
		resultRecord, err := ScanRowToRecord(row)

		if !(testRecord[0].Equals(resultRecord)) {
			t.Fatalf("Record mismatch.\n\nExpected:\n\n%v\n\nReceived:\n\n%v", testRecord[0], resultRecord)
		}
	})

	t.Run("inserts multiple records correctly", func(t *testing.T) {
		setup, err := Setup()
		if err != nil {
			t.Fatalf("Failed to setup test suite: %v", err)
		}
		defer setup.TearDown()

		testRecords := []mds.MarketHistoryCSVRecord{
			{
				Date:       time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
				RegionID:   1,
				TypeID:     1,
				Average:    1,
				Highest:    2,
				Lowest:     1,
				Volume:     1,
				OrderCount: 1,
			},
			{
				Date:       time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
				RegionID:   1,
				TypeID:     2,
				Average:    1,
				Highest:    2,
				Lowest:     1,
				Volume:     1,
				OrderCount: 1,
			},
		}

		if err = (*setup.DBService).InsertMany(testRecords); err != nil {
			t.Fatalf("Failed to insert multiple records: %v", err)
		}

		rows, err := setup.Connection.Query("SELECT * FROM market_data")
		if err != nil {
			t.Fatalf("Failed to retrieve inserted Rows for comparison: %v", err)
		}

		results := make([]mds.MarketHistoryCSVRecord, 0, 2)
		for rows.Next() {
			scannedRecord, err := ScanRowToRecord(rows)
			if err != nil {
				t.Fatalf("Failed to scan record: %v", err)
			}
			results = append(results, *scannedRecord)
		}

		if len(results) != 2 {
			t.Fatalf("Expected 2 results, received: %d", len(results))
		}

		for i := range len(results) {
			if !testRecords[i].Equals(&results[i]) {
				t.Fatalf("Record mismatch.\n\nExpected:\n\n%v\n\nReceived:\n\n%v", testRecords[i], results[i])
			}
		}
	})

	t.Run("inserts an entire data file correctly", func(t *testing.T) {
		// arrange
		setup, err := Setup()
		if err != nil {
			t.Fatalf("Failed to setup test suite: %v", err)
		}
		defer setup.TearDown()

		// get a marketdataservice
		mockDownload := func(_ string) (mds.ZippedReader, error) {
			file, err := os.ReadFile("../testdata/exports/market-history-2024-02-10.csv.bz2")
			if err != nil {
				t.Fatalf("Failed to open csv file: %v", err)
			}
			return bytes.NewReader(file), nil
		}

		dataService := mds.NewMarketDataService(mockDownload, mds.DecompressFile, mds.ParseFile)
		records, err := dataService.FetchAndParseCSV("website.ext/file.ext")
		if err != nil {
			t.Fatalf("Failed to parse records from file: %v", err)
		}

		// act
		if err = (*setup.DBService).InsertMany(records); err != nil {
			t.Fatalf("Failed to insert records: %v", err)
		}

		// assert
		row := setup.Connection.QueryRow("SELECT COUNT(date) FROM market_data")
		if err != nil {
			t.Fatalf("Failed to query number of rows")
		}
		var numRows int
		row.Scan(&numRows)

		if len(records) != numRows {
			t.Fatalf("Expected to see %d records, actually saw %d.", len(records), numRows)
		}
	})
}
