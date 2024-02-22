package marketdataservice

import (
	"io"
	"time"
)

type MarketHistoryCSVRecord struct {
	Date       time.Time `db:"date"`
	RegionID   uint      `db:"region_id"`
	TypeID     uint      `db:"type_id"`
	Average    float64   `db:"average"`
	Highest    float64   `db:"highest"`
	Lowest     float64   `db:"lowest"`
	Volume     uint      `db:"volume"`
	OrderCount uint      `db:"order_count"`
}

type MarketDay struct {
	Date    time.Time
	Records []MarketHistoryCSVRecord
}

// semantic wrappers for explicit format safety
type ZippedReader io.Reader
type UnzippedReader io.Reader

type Download func(url string) (ZippedReader, error)
type Decompress func(reader ZippedReader) (UnzippedReader, error)
type Parse func(reader UnzippedReader) ([]MarketHistoryCSVRecord, error)

type IMarketDataService interface {
	FetchAndParseDay(day time.Time) (*MarketDay, error)
}

type MarketDataService struct {
	download   Download
	decompress Decompress
	parse      Parse
}
