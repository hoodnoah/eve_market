package integration_tests

import (
	"compress/bzip2"
	"os"
	"testing"
	"time"

	"github.com/hoodnoah/eve_market/monitor/parser"
)

func outputChanToSlice(outputChan chan time.Time) []time.Time {
	output := make([]time.Time, 0, 1)
	timeout := time.After(2000 * time.Millisecond)
	for {
		select {
		case d := <-outputChan:
			output = append(output, d)
		case <-timeout:
			return output
		}
	}
}

func TestInsertMarketDay(t *testing.T) {
	t.Run("it should insert a sample day with one record", func(t *testing.T) {
		// arrange
		inputChan := make(chan *parser.MarketDay)
		outputChan := make(chan time.Time)

		testSetup, err := Setup(inputChan, outputChan)
		if err != nil {
			t.Fatalf("failed to get test setup: %v", err)
		}
		testSetup.DBManager.Start()
		defer testSetup.TearDown()

		testRecords := []parser.MarketHistoryCSVRecord{
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

		testDay := parser.MarketDay{
			Date:    time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC),
			Records: testRecords,
		}

		// act
		inputChan <- &testDay
		outputs := outputChanToSlice(outputChan)

		if len(outputs) != 1 {
			t.Fatalf("Expected 1 output for 1 day's insertion, received %d", len(outputs))
		}

		if !outputs[0].Equal(time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("Expected returned date to be 2003-10-01, received %s", outputs[0].Format(time.DateOnly))
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

	t.Run("it should insert the sample file from 2024-02-10", func(t *testing.T) {
		// arrange
		inputChan := make(chan *parser.MarketDay)
		outputChan := make(chan time.Time)
		testSetup, err := Setup(inputChan, outputChan)
		testSetup.DBManager.Start()
		if err != nil {
			t.Fatalf("failed to get test setup: %v", err)
		}
		defer testSetup.TearDown()

		// read in the input file and parse it
		inputFile, err := os.Open("../testdata/exports/market-history-2024-02-10.csv.bz2")
		if err != nil {
			t.Fatalf("failed to read input file: %v", err)
		}
		reader := parser.DatedReader{
			Day:    time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			Reader: parser.UnzippedReader(bzip2.NewReader(inputFile)),
		}
		decompressed, err := parser.ParseFile(&reader)
		if err != nil {
			t.Fatalf("failed to parse the input file: %v", err)
		}

		// act
		inputChan <- decompressed
		outputList := outputChanToSlice(outputChan)

		if len(outputList) != 1 {
			t.Fatalf("expected only 1 returned date, received %d", len(outputList))
		}

		if !outputList[0].Equal(time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("expected to receive 2024-02-10, received %s", outputList[0].Format(time.DateOnly))
		}

		// assert
		// should only insert 1 day into completed_dates
		var numDates uint
		if err = testSetup.Connection.QueryRow("SELECT COUNT(id) FROM completed_dates;").Scan(&numDates); err != nil {
			t.Fatalf("failed to retrieve row count: %v", err)
		}
		if numDates != 1 {
			t.Fatalf("Expected 1 date, counted %d", numDates)
		}

		// should insert the full 52,484 records
		var numRecords uint
		if err = testSetup.Connection.QueryRow("SELECT COUNT(date_id) FROM market_data").Scan(&numRecords); err != nil {
			t.Fatalf("failed to retrieve row count from market_data: %v", numRecords)
		}

		if numRecords != 52484 {
			t.Fatalf("Expoected 52484 records, received %d", numRecords)
		}
	})
}
