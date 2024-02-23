package integration_tests

import (
	"testing"
	"time"

	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

func TestInsertMarketDay(t *testing.T) {
	t.Run("it should insert a sample day with one record", func(t *testing.T) {
		// arrange
		testSetup, err := Setup()
		if err != nil {
			t.Fatalf("failed to get test setup: %v", err)
		}
		// defer testSetup.TearDown()

		testRecords := []mds.MarketHistoryCSVRecord{
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

		testDay := mds.MarketDay{
			Date:    time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC),
			Records: testRecords,
		}

		// act
		_, err = testSetup.DBManager.InsertMarketDay(&testDay)
		if err != nil {
			t.Fatalf("failed to insert a day: %v", err)
		}

		// assert
		// should have inserted a new date
		var id int
		var date time.Time
		if err = testSetup.Connection.QueryRow("SELECT id, date FROM completed_dates").Scan(&id, &date); err != nil {
			t.Fatalf("failed to query a row from completed_dates: %v", err)
		}

		if id != 1 {
			t.Fatalf("Expected an id of 1, received: %v", id)
		}
		if !date.Equal(time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("Expected to date 2003-10-01, received: %s", date.String())
		}

		// should have inserted a new record
		var dateId uint
		var regionId uint
		var typeId uint
		var average float64
		var highest float64
		var lowest float64
		var volume uint
		var orderCount uint
		if err = testSetup.Connection.QueryRow("SELECT * FROM market_data").Scan(
			&dateId,
			&regionId,
			&typeId,
			&average,
			&highest,
			&lowest,
			&volume,
			&orderCount,
		); err != nil {
			t.Fatalf("Failed to fetch record from market_data: %v", err)
		}

		tr := testRecords[0]
		switch {
		case dateId != 1:
			t.Fatalf("expected dateID %d, received %d", 1, dateId)
		case tr.RegionID != regionId:
			t.Fatalf("expected regionID %d, received %d", tr.RegionID, regionId)
		case tr.TypeID != typeId:
			t.Fatalf("expected typeID %d, received %d", tr.TypeID, typeId)
		case tr.Average != average:
			t.Fatalf("expected average %f, received %f", tr.Average, average)
		case tr.Highest != highest:
			t.Fatalf("expected highest %f, received %f", tr.Highest, highest)
		case tr.Lowest != lowest:
			t.Fatalf("expected lowest %f, received %f", tr.Lowest, lowest)
		case tr.Volume != volume:
			t.Fatalf("expected volume %d, received %d", tr.Volume, volume)
		case tr.OrderCount != orderCount:
			t.Fatalf("expected orderCount %d, received %d", tr.OrderCount, orderCount)
		}

		// should only have inserted a single record into either table
		var numRecords int
		if err = testSetup.Connection.QueryRow("SELECT COUNT(date_id) FROM market_data").Scan(&numRecords); err != nil {
			t.Fatalf("failed to count the records in market_data: %v", err)
		}
		if numRecords != 1 {
			t.Fatalf("expected 1 record in market_data, received %d", numRecords)
		}

		if err = testSetup.Connection.QueryRow("SELECT COUNT(id) FROM completed_dates").Scan(&numRecords); err != nil {
			t.Fatalf("failed to count the records in completed_dates: %v", err)
		}
		if numRecords != 1 {
			t.Fatalf("expected 1 record in completed_dates, received %d", numRecords)
		}

	})
}
