package dates

import (
	"sync"
	"time"
)

var firstDateWithData = time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)

func NewDateIterator(currentDateFn NowFn) *DateIterator {
	return &DateIterator{
		currentDate: firstDateWithData,
		nowFn:       currentDateFn,
		mutex:       sync.Mutex{},
	}
}

// gets the next date in sequence
// 1 day after the last date, starting from 2003-10-01
func (di *DateIterator) Next() time.Time {
	di.mutex.Lock()
	defer di.mutex.Unlock()

	returnValue := di.currentDate
	di.currentDate = di.currentDate.AddDate(0, 0, 1)
	return returnValue
}

// returns true if the next date in the sequence
// is, at the latest, yesterday
func (di *DateIterator) IsNextReady() bool {
	year, month, day := di.nowFn().Date()
	cutoffDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	return cutoffDate.After(di.currentDate)
}

// tautology; it's an infinite
func (di *DateIterator) HasNext() bool {
	return true
}
