package marketdataservice

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

func linearSearch(list []string, target string) int {
	for i, elem := range list {
		if strings.ToLower(elem) == strings.ToLower(target) {
			return i
		}
	}
	return -1
}

func findColumnIndices(headerRecord []string) (map[string]int, error) {
	if len(headerRecord) < 8 {
		return nil, fmt.Errorf("header row too short; expected 8 columns minimum, received %d", len(headerRecord))
	}

	columnIndices := make(map[string]int)
	for _, columnHeader := range fieldNames {
		index := linearSearch(headerRecord, columnHeader)
		if index < 0 {
			return nil, fmt.Errorf("header row missing required field %s", columnHeader)
		}
		if index > 0 {
			columnIndices[columnHeader] = index
		}
	}

	return columnIndices, nil
}

func parseRecord(columnIndices map[string]int, record []string) (*MarketHistoryCSVRecord, error) {
	parsedDate, err := parseDate(record[columnIndices["date"]])

	if err != nil {
		return nil, err
	}

	parsedRegionID, err := parseUint(record[columnIndices["region_id"]])

	if err != nil {
		return nil, err
	}

	parsedTypeID, err := parseUint(record[columnIndices["type_id"]])

	if err != nil {
		return nil, err
	}

	parsedAverage, err := parseFloat(record[columnIndices["average"]])

	if err != nil {
		return nil, err
	}

	parsedHighest, err := parseFloat(record[columnIndices["highest"]])

	if err != nil {
		return nil, err
	}

	parsedLowest, err := parseFloat(record[columnIndices["lowest"]])

	if err != nil {
		return nil, err
	}

	parsedVolume, err := parseUint(record[columnIndices["volume"]])

	if err != nil {
		return nil, err
	}

	parsedOrderCount, err := parseUint(record[columnIndices["order_count"]])

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

// converts a date to its url stub
// e.g. "2003-10-01" -> ".../2003/market-history-2003-10-01.csv.bz2"
func dateToURL(date time.Time) string {
	return fmt.Sprintf(
		urlTemplate,
		date.Year(),
		date.Year(),
		date.Month(),
		date.Day(),
	)
}
