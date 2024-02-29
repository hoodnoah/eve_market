package util

// chunk a slice into slices of at most size n, where
// n is the provided chunkSize
func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	returnValue := make([][]T, 0)
	if len(slice) <= chunkSize {
		return append(returnValue, slice)
	}

	idx := 0
	for idx < len(slice) {
		high := idx + min(chunkSize, len(slice)-idx)
		returnValue = append(returnValue, slice[idx:high])
		idx = high
	}

	return returnValue
}
