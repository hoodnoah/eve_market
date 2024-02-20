package marketdataservice

import (
	"compress/bzip2"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var fieldNames = [8]string{
	"date",
	"region_id",
	"type_id",
	"average",
	"highest",
	"lowest",
	"volume",
	"order_count",
}

// parses a date of format yyyy-mm-dd into a time
func parseDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, dateStr)

	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func parseUint(uintStr string) (uint, error) {
	result, err := strconv.ParseUint(uintStr, 10, 64)

	if err != nil {
		return 0, err
	}

	return uint(result), nil
}

func parseFloat(floatStr string) (float64, error) {
	result, err := strconv.ParseFloat(floatStr, 64)

	if err != nil {
		return 0, err
	}

	return result, nil
}

// Determines if a slice contains the correct header field names
func isCorrectHeaderRow(record []string) bool {
	if len(record) != 8 {
		return false
	}
	for i, val := range record {
		if val != fieldNames[i] {
			return false
		}
	}

	return true
}

func parseRecord(record []string) (*MarketHistoryCSVRecord, error) {
	// expect date, region_id, type_id, average, highest, lowest, volume, order_count
	if len(record) != 8 {
		return nil, fmt.Errorf("expected 8 fields, received %d", len(record))
	}

	var err error

	parsedDate, err := parseDate(record[0])

	if err != nil {
		return nil, err
	}

	parsedRegionID, err := parseUint(record[1])

	if err != nil {
		return nil, err
	}

	parsedTypeID, err := parseUint(record[2])

	if err != nil {
		return nil, err
	}

	parsedAverage, err := parseFloat(record[3])

	if err != nil {
		return nil, err
	}

	parsedHighest, err := parseFloat(record[4])

	if err != nil {
		return nil, err
	}

	parsedLowest, err := parseFloat(record[5])

	if err != nil {
		return nil, err
	}

	parsedVolume, err := parseUint(record[6])

	if err != nil {
		return nil, err
	}

	parsedOrderCount, err := parseUint(record[7])

	if err != nil {
		return nil, err
	}

	return &MarketHistoryCSVRecord{
		Date:       parsedDate,
		RegionID:   parsedRegionID,
		TypeID:     parsedTypeID,
		Average:    parsedAverage,
		Highest:    parsedHighest,
		Lowest:     parsedLowest,
		Volume:     parsedVolume,
		OrderCount: parsedOrderCount,
	}, nil

}

// constructor for marketdataservice struct
func NewMarketDataService(download Download, decompress Decompress, parse Parse) *MarketDataService {
	return &MarketDataService{
		download:   download,
		decompress: decompress,
		parse:      parse,
	}
}

// Downloads the file at the provided URL,
// returning the body Reader
func DownloadFile(url string) (ZippedReader, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("failed to retrieve the file at %s with error code %d: %s", url, res.StatusCode, res.Status)
	}
	return res.Body, nil
}

// Decompresses the provided bzip file
func DecompressFile(zippedFile ZippedReader) (UnzippedReader, error) {
	return bzip2.NewReader(zippedFile), nil
}

// Parses the provided file into a slice of MarketHistoryCSVRecords
func ParseFile(unzippedFile UnzippedReader) ([]MarketHistoryCSVRecord, error) {
	// step 1: read into records
	r := csv.NewReader(unzippedFile)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if !isCorrectHeaderRow(records[0]) {
		return nil, errors.New("received malformed header row")
	}

	var returnRecords []MarketHistoryCSVRecord

	for _, val := range records[1:] {
		parseResult, err := parseRecord(val)
		if err != nil {
			return nil, err
		}
		returnRecords = append(returnRecords, *parseResult)
	}

	return returnRecords, nil
}

// Fetches (by URL) and parses a MarketHistoryCSV file into its constituent records
func (m *MarketDataService) FetchAndParseCSV(url string) ([]MarketHistoryCSVRecord, error) {
	// Download
	bzippedData, err := m.download(url)
	if err != nil {
		return nil, err
	}

	// Decompress
	unzippedData, err := m.decompress(bzippedData)
	if err != nil {
		return nil, err
	}

	// Parse
	records, err := m.parse(unzippedData)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// compares the current MarketDataCSVRecord with a provided one
func (m *MarketHistoryCSVRecord) Equals(other *MarketHistoryCSVRecord) bool {
	switch {
	case !m.Date.Equal(other.Date):
		return false
	case m.RegionID != other.RegionID:
		return false
	case m.TypeID != other.TypeID:
		return false
	case m.Average != other.Average:
		return false
	case m.Highest != other.Highest:
		return false
	case m.Lowest != other.Lowest:
		return false
	case m.Volume != other.Volume:
		return false
	case m.OrderCount != other.OrderCount:
		return false
	default:
		return true
	}
}
