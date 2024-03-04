package parser

// compares the current MarketDataCSVRecord with a provided one
func (m *MarketHistoryCSVRecord) Equals(other *MarketHistoryCSVRecord) bool {
	switch {
	case !m.Date.Equal(other.Date):
		return false
	case m.RegionID != other.RegionID:
		return false
	case m.TypeID != other.TypeID:
		return false
	case m.Average != other.Average:
		return false
	case m.Highest != other.Highest:
		return false
	case m.Lowest != other.Lowest:
		return false
	case m.Volume != other.Volume:
		return false
	case m.OrderCount != other.OrderCount:
		return false
	default:
		return true
	}
}
