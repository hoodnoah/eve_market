package timeutils

import (
	"fmt"
	"time"
)

// earliest year with data
const EarliestYear = 2003

// the prefix for the daily export data file
const FilePrefix = "market-history"

// the suffix for the daily export data file
const FileExt = ".csv.bz2"

// first date with data
var EarliestDate = time.Date(EarliestYear, 10, 1, 0, 0, 0, 0, time.UTC)

// enumerates the years from the beginning of daily market data collection (2003)
// to present
func EnumerateYearsToPresent(startYear int) []int {
	currentYear := time.Now().Year()
	numYears := int(currentYear) - int(startYear) + 1
	years := make([]int, numYears)

	for i := range years {
		years[i] = int(startYear) + i
	}

	return years
}

// creates a slice containing all the days within a given integer year
// excluding the beginning of 2003, since data recording began on
// 2003-10-01
func EnumerateDatesInYear(year int) []time.Time {
	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)

	var dates []time.Time

	for currentDate := startDate; currentDate.Before(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if currentDate.After(EarliestDate) || currentDate.Equal(EarliestDate) {
			dates = append(dates, currentDate)
		}
	}

	return dates
}

// generates a url suffix for a given day's data,
// e.g. for '2003-10-01' -> 2003/market-history-2003-10-01.csv.bz2'
func GetDayFileURL(date time.Time) string {
	return fmt.Sprintf("%d/%s-%s%s", date.Year(), FilePrefix, date.Format(time.DateOnly), FileExt)
}
