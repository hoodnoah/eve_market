package integration_tests

import (
	"testing"
	"time"

	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

func TestInsertMany(t *testing.T) {
	t.Run("Inserts a single sample record correctly", func(t *testing.T) {
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
		var resultRecord mds.MarketHistoryCSVRecord
		row := setup.Connection.QueryRow(`
			SELECT * FROM market_data;
			`)
		row.Scan(
			&resultRecord.Date,
			&resultRecord.RegionID,
			&resultRecord.TypeID,
			&resultRecord.Average,
			&resultRecord.Highest,
			&resultRecord.Lowest,
			&resultRecord.Volume,
			&resultRecord.OrderCount,
		)

		if !(testRecord[0].Equals(&resultRecord)) {
			t.Fatalf("Record mismatch.\n\nExpected:\n\n%v\n\nReceived:\n\n%v", testRecord[0], resultRecord)
		}
	})
}
