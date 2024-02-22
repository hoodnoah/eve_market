package marketdataservice

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func recordEqual(r1 MarketHistoryCSVRecord, r2 MarketHistoryCSVRecord) bool {
	switch {
	case !r1.Date.Equal(r2.Date):
		return false
	case r1.RegionID != r2.RegionID:
		return false
	case r1.TypeID != r2.TypeID:
		return false
	case r1.Average != r2.Average:
		return false
	case r1.Highest != r2.Highest:
		return false
	case r1.Lowest != r2.Lowest:
		return false
	case r1.Volume != r2.Volume:
		return false
	case r1.OrderCount != r2.OrderCount:
		return false
	default:
		return true
	}
}

func slicesEqual(s1 []MarketHistoryCSVRecord, s2 []MarketHistoryCSVRecord, comp func(MarketHistoryCSVRecord, MarketHistoryCSVRecord) bool) bool {
	if len(s1) != len(s2) {
		return false
	}

	// lengths are the same
	for i := range s1 {
		if !comp(s1[i], s2[i]) {
			return false
		}
	}
	return true
}

const file2003 = "./testdata/marketdataservice/market-history-2003-10-01.csv.bz2"
const file2024 = "./testdata/marketdataservice/market-history-2024-02-10.csv.bz2"

var EXPECTED_RECORDS = []MarketHistoryCSVRecord{
	MarketHistoryCSVRecord{
		Date:       time.Date(2003, time.October, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     18,
		Average:    10,
		Highest:    10,
		Lowest:     10,
		Volume:     340,
		OrderCount: 2,
	},
	MarketHistoryCSVRecord{
		Date:       time.Date(2003, time.October, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     20,
		Average:    14,
		Highest:    14,
		Lowest:     14,
		Volume:     143,
		OrderCount: 1,
	},
	MarketHistoryCSVRecord{
		Date:       time.Date(2003, time.October, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     34,
		Average:    1,
		Highest:    1,
		Lowest:     1,
		Volume:     21188886,
		OrderCount: 43,
	},
	MarketHistoryCSVRecord{
		Date:       time.Date(2003, time.October, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     35,
		Average:    4,
		Highest:    4,
		Lowest:     4,
		Volume:     9250727,
		OrderCount: 40,
	},
	MarketHistoryCSVRecord{
		Date:       time.Date(2003, time.October, 1, 0, 0, 0, 0, time.UTC),
		RegionID:   10000001,
		TypeID:     36,
		Average:    16,
		Highest:    16,
		Lowest:     16,
		Volume:     216188,
		OrderCount: 20,
	},
}

func setup(filePath string) struct {
	Service IMarketDataService
} {
	mockDownload := func(s string) (ZippedReader, error) {
		file, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(file), nil
	}

	return struct {
		Service IMarketDataService
	}{
		Service: NewMarketDataService(mockDownload, DecompressFile, ParseFile),
	}
}

func TestMarketDataService_FetchAndParseCSV(t *testing.T) {
	t.Run("reads truncated 2003 file correctly", func(t *testing.T) {
		fixture := setup(file2003)
		service := fixture.Service

		testDate := time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)

		day, err := service.FetchAndParseDay(testDate)
		if err != nil {
			t.Fatalf("FetchAndParseDay returned an error: %v", err)
		}

		if !day.Date.Equal(testDate) {
			t.Fatalf("Returned incorrect date. Expected %s, received %s", testDate.String(), day.Date.String())
		}

		if !slicesEqual(day.Records, EXPECTED_RECORDS, recordEqual) {
			t.Fatalf("Expected:\n%v\nReceived:\n%v", EXPECTED_RECORDS, day.Records)
		}
	})

	t.Run("reads full 2024 file correctly", func(t *testing.T) {
		fixture := setup(file2024)

		testDate := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
		expectedNumRecords := 52484

		day, err := fixture.Service.FetchAndParseDay(testDate)
		if err != nil {
			t.Fatalf("FetchAndParseDay returned an error: %v", err)
		}

		if !day.Date.Equal(testDate) {
			t.Fatalf("Returned incorrect date. Expected %s, received %s", testDate.String(), day.Date.String())
		}

		if len(day.Records) != expectedNumRecords {
			t.Fatalf("Expected to receive %d records, received %d", expectedNumRecords, len(day.Records))
		}
	})
}
