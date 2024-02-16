package datadateservice

import "time"

type Year int

// a single date
type DataDate struct {
	Date time.Time
	Url  string
}

// a year's worth of dates
type DataYear struct {
	Year  Year
	Dates []DataDate
}

type IDataDateService interface {
	EnumerateDataYears(currentDate time.Time) []DataYear
}

type DataDateService struct {
	getCurrentDate func() time.Time
}
