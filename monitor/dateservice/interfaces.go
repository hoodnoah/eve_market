package dateservice

import "time"

type Year int

type IDateService interface {
	// generates a list of dates
	EnumerateDates() []time.Time
}

type EveRefsDateService struct {
	datefn func() time.Time
}
