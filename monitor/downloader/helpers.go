package downloader

import (
	"fmt"
	"time"
)

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
