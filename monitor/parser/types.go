package parser

import (
	"io"
	"sync"
	"time"

	"github.com/hoodnoah/eve_market/monitor/logger"
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

type DatedReader struct {
	Day    time.Time
	Reader io.Reader
}

type Decompress func(zippedRecord *DatedReader) *DatedReader
type Parse func(unzippedRecord *DatedReader) (*MarketDay, error)

type IMarketDataParser interface {
	Start()
}

type MarketDataParser struct {
	logger     logger.ILogger
	decompress Decompress
	parse      Parse
	numWorkers uint
	input      chan *DatedReader
	results    chan *MarketDay
	mutex      sync.Mutex
}
