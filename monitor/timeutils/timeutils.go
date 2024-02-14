package timeutils

import (
	"fmt"
	"time"
)

// semantic type alias
type Year uint

// earliest year with data
const earliestYear = Year(2003)

// the prefix for the daily export data file
const filePrefix = "market-history"

// the suffix for the daily export data file
const fileExt = ".csv.bz2"

// first date with data
var earliestDate = time.Date(int(earliestYear), 10, 1, 0, 0, 0, 0, time.UTC)

// enumerates the years from the beginning of daily market data collection (2003)
// to present
func EnumerateYearsToPresent(startYear Year) []Year {
	currentYear := Year(time.Now().Year())
	numYears := currentYear - earliestYear + 1
	years := make([]Year, numYears)

	for i := range years {
		years[i] = startYear + Year(i)
	}

	return years
}

// creates a slice containing all the days within a given integer year
// excluding the beginning of 2003, since data recording began on
// 2003-10-01
func EnumerateDatesInYear(year Year) []time.Time {
	startDate := time.Date(int(year), time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(int(year+1), time.January, 1, 0, 0, 0, 0, time.UTC)

	var dates []time.Time

	for currentDate := startDate; currentDate.Before(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if currentDate.After(earliestDate) || currentDate.Equal(earliestDate) {
			dates = append(dates, currentDate)
		}
	}

	return dates
}

// generates a url suffix for a given day's data,
// e.g. for '2003-10-01' -> 2003/market-history-2003-10-01.csv.bz2'
func GetDayFileURL(date time.Time) string {
	return fmt.Sprintf("%d/%s-%s%s", date.Year(), filePrefix, date.Format(time.DateOnly), fileExt)
}
