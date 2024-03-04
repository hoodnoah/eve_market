package idcache

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/util"
)

func NewIDCache(logger logger.ILogger, client IRequestClient) IIDcache {
	return &IDCache{
		logger: logger,
		client: client,
		ids:    map[int]string{},
		mutex:  sync.Mutex{},
	}
}

func (cache *IDCache) Label(id int) (string, error) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// if cache doesn't contain the id,
	// fetch it
	if cache.ids[id] == "" {
		val, err := cache.client.FetchID(id)
		if err != nil {
			switch err.ErrorType {
			case InvalidIDRequestError: // API says the ID doesn't exist, so use a default identifier
				cache.ids[id] = fmt.Sprintf("invalidID_%d", id)
			default: // otherwise do nothing; other errors are undefined behavior
			}
		} else {
			cache.ids[id] = val.Name // success case, add to the cache
		}
	}

	// if the cache contains the id, return it
	if cache.ids[id] != "" {
		return cache.ids[id], nil
	}

	return "", errors.New("ID unresolvable")
}

// fetches a series of ids.
// if, at any point, not all ids can be fetched and the reason is determined
// to be a 404 error, indicating an invalid ID, it recursively narrows the search space until
// the individual invalid ID can be detected and assigned an unknownID_id value.
// if it fails for any other reason, the entire operation fails.
func (idc *IDCache) fetchAllIds(ids []int) ([]ESITypeIDResponse, *APIRequestError) {
	// base case 1: only 1 id left
	if len(ids) == 1 {
		id := ids[0]
		res, err := idc.client.FetchID(id)
		if err != nil && err.ErrorType == InvalidIDRequestError { // invalid id detected
			result := ESITypeIDResponse{
				Category: "",
				ID:       id,
				Name:     fmt.Sprintf("unknownID_%d", id),
			}
			return []ESITypeIDResponse{result}, nil
		}
		if err != nil && err.ErrorType != InvalidIDRequestError { // some other error
			return nil, err
		}
		return []ESITypeIDResponse{*res}, nil
	}

	// base case 2: all ids resolve correctly, or error
	// for a reason other than invalidIDRequestError
	allIdsRes, allIdsErr := idc.client.FetchManyIDs(ids)
	if allIdsErr == nil {
		return allIdsRes, nil
	}
	if allIdsErr.ErrorType != InvalidIDRequestError {
		return nil, allIdsErr
	}

	// recurse
	midPoint := len(ids) / 2
	left := ids[0:midPoint]
	right := ids[midPoint:]

	result1, err1 := idc.fetchAllIds(left)
	result2, err2 := idc.fetchAllIds(right)

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	return append(result1, result2...), nil
}

func filterToUniqueIds(ids []int) []int {
	idsMap := map[int]bool{}
	output := make([]int, 0)
	for _, id := range ids {
		idsMap[id] = true
	}

	for id, val := range idsMap {
		if val {
			output = append(output, id)
		}
	}

	return output
}

// label many ids at once
func (idc *IDCache) LabelMany(ids []int) (map[int]string, error) {
	uniqueIds := filterToUniqueIds(ids)
	result := map[int]string{}

	// find unknown ids
	novelIds := idc.FindNovelIDs(uniqueIds)
	if len(novelIds) == 0 {
		for _, id := range uniqueIds {
			result[id] = idc.ids[id]
		}
	} else {
		chunks := util.ChunkSlice(novelIds, 1000)
		for _, chunk := range chunks {
			// fetch unknown ids
			esiResults, err := idc.fetchAllIds(chunk)

			if err != nil {
				return nil, err
			}

			// add to cache
			for _, item := range esiResults {
				idc.ids[item.ID] = item.Name
			}
		}
	}

	// populate response
	for _, id := range uniqueIds {
		result[id] = idc.ids[id]
	}

	return result, nil

}

func (idc *IDCache) FindNovelIDs(ids []int) []int {
	filteredIDs := make([]int, 0)
	for _, id := range ids {
		if idc.ids[id] == "" {
			filteredIDs = append(filteredIDs, id)
		}
	}

	return filteredIDs
}

// sets the ids already known
// should generally be retrieved from the backing db,
// then sent here
func (idc *IDCache) SetKnownIDs(ids map[int]string) {
	idc.mutex.Lock()
	defer idc.mutex.Unlock()

	for id, val := range ids {
		idc.ids[id] = val
	}
}
