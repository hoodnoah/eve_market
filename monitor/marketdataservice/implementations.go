package marketdataservice

import (
	"compress/bzip2"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"
)

// compile-time interface check
var _ IMarketDataService = (*MarketDataService)(nil)

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

	columnIndices, err := findColumnIndices(records[0])
	if err != nil {
		return nil, err
	}

	var returnRecords []MarketHistoryCSVRecord

	for _, val := range records[1:] {
		parseResult, err := parseRecord(columnIndices, val)
		if err != nil {
			return nil, err
		}
		returnRecords = append(returnRecords, *parseResult)
	}

	return returnRecords, nil
}

// Fetches (by URL) and parses a MarketHistoryCSV file into its constituent records
func (m *MarketDataService) FetchAndParseDay(day time.Time) (*MarketDay, error) {
	// Download
	bzippedData, err := m.download(dateToURL(day))
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

	return &MarketDay{
		Date:    day,
		Records: records,
	}, nil
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
