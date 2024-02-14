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

const ZIPPED_CSV_PATH = "./testdata/marketdataservice/compressed.csv.bz2"

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

// simulates the downloading of a bzipped csv file
// by reading in a test fixture
func mockDownload(_ string) (ZippedReader, error) {
	file, err := os.ReadFile(ZIPPED_CSV_PATH)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(file), nil
}

func TestMarketDataService_FetchAndParseCSV(t *testing.T) {
	service := NewMarketDataService(mockDownload, DecompressFile, ParseFile)

	records, err := service.FetchAndParseCSV("website.ext/file.ext")
	if err != nil {
		t.Fatalf("FetchAndParseCSV returned an error: %v", err)
	}

	if !slicesEqual(records, EXPECTED_RECORDS, recordEqual) {
		t.Fatalf("Expected:\n%v\nReceived:\n%v", EXPECTED_RECORDS, records)
	}
}
