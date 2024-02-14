package marketdataservice

import (
	"io"
	"time"
)

type MarketHistoryCSVRecord struct {
	Date       time.Time
	RegionID   uint
	TypeID     uint
	Average    float64
	Highest    float64
	Lowest     float64
	Volume     uint
	OrderCount uint
}

// semantic wrappers for explicit format safety
type ZippedReader io.Reader
type UnzippedReader io.Reader

type Download func(url string) (ZippedReader, error)
type Decompress func(reader ZippedReader) (UnzippedReader, error)
type Parse func(reader UnzippedReader) ([]MarketHistoryCSVRecord, error)

type IMarketDataService interface {
	FetchAndParseCSV(url string) ([]MarketHistoryCSVRecord, error)
}

type MarketDataService struct {
	download   Download
	decompress Decompress
	parse      Parse
}
