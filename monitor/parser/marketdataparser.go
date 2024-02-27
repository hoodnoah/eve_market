package parser

import (
	"fmt"
	"time"

	"github.com/hoodnoah/eve_market/monitor/logger"
)

// compile-time interface check
var _ IMarketDataParser = (*MarketDataParser)(nil)

// constructor for marketdataservice struct
func NewMarketDataParser(logger logger.ILogger, decompress Decompress, parse Parse, inputChannel chan *DatedReader, resultsChannel chan *MarketDay, numWorkers uint) IMarketDataParser {
	if logger == nil {
		panic("Failed to create MarketDataParser: logger cannot be nil")
	}

	return &MarketDataParser{
		decompress: decompress,
		parse:      parse,
		logger:     logger,
		input:      inputChannel,
		results:    resultsChannel,
		numWorkers: numWorkers,
	}
}

func (m *MarketDataParser) decompressAndParse(zippedFile *DatedReader) {
	decompressed := m.decompress(zippedFile)
	parsed, err := m.parse(decompressed)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse the file for %s, %v", zippedFile.Day.Format(time.DateOnly), err)
		m.logger.Error(errMsg)
	} else {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.results <- &MarketDay{
			Date:    parsed.Date,
			Records: parsed.Records,
		}
	}
}

func (m *MarketDataParser) Start() {
	for range m.numWorkers {
		go func() {
			for reader := range m.input {
				m.decompressAndParse(reader)
			}
		}()
	}
}
