package main

import (
	dateservice "github.com/hoodnoah/eve_market/monitor/datadateservice"
)

func main() {
	// setup services
	dateSvc := dateservice.NewDataDateService(dateservice.GetCurrentDate)

	// get dates
	dataYears := dateSvc.EnumerateDataYears()
}
