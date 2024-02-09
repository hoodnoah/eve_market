package main

import (
	"fmt"

	"github.com/hoodnoah/eve_market/monitor/timeutils"
)

const BaseUrl = "https://data.everef.net/market-history/"

func main() {
	// Get list of years from the first time data was listed to present
	// yearsList := timeutils.EnumerateYearsToPresent(timeutils.EarliestYear)

	// Get days in 2003
	days2003 := timeutils.EnumerateDatesInYear(2003)

	firstDay := days2003[0]

	fmt.Println(firstDay)
}
