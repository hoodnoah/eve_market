package idcache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/ratelimiter"
	"github.com/hoodnoah/eve_market/monitor/util"
)

func NewIDManager(logger logger.ILogger) IIDCache {
	tokenBucket := ratelimiter.NewTokenBucketRateLimiter(1)
	tokenBucket.Start()

	return &IDCache{
		logger:      logger,
		regionIDS:   map[int]string{},
		typeIDS:     map[int]string{},
		rateLimiter: tokenBucket,
		client:      *getClient(),
		mutex:       sync.Mutex{},
	}
}

func (idmgr *IDCache) SetKnownRegionIDs(regionIds *RegionIDInput) {
	idmgr.mutex.Lock()
	defer idmgr.mutex.Unlock()

	idmgr.regionIDS = regionIds.IDS
}

func (idmgr *IDCache) SetKnownTypeIDs(typeIds *TypeIDInput) {
	idmgr.mutex.Lock()
	defer idmgr.mutex.Unlock()

	idmgr.typeIDS = typeIds.IDS
}

// given unlabeled ids, map labels to them.
// fetches their labels from the EVE API as needed
func (idm *IDCache) Label(ids *UnknownIDs) (*KnownIDs, error) {
	idm.mutex.Lock()
	defer idm.mutex.Unlock()

	// set map of known ids based on type
	var knownIds *map[int]string
	switch ids.Type {
	case RegionID:
		knownIds = &idm.regionIDS
	default:
		knownIds = &idm.typeIDS
	}

	idsToFetch := map[int]bool{}
	for id := range ids.IDS {
		if (*knownIds)[id] == "" {
			idsToFetch[id] = true
		}
	}

	newIds, err := idm.fetchUnknownIds(&idsToFetch)
	if err != nil {
		idm.logger.Error(fmt.Sprintf("failed to fetch unknown ids: %s", err))
		return nil, err
	}

	for id, value := range newIds {
		(*knownIds)[id] = value
	}

	output := KnownIDs{
		Type: ids.Type,
		IDS:  *knownIds,
	}

	return &output, nil
}

func (idm *IDCache) fetchUnknownIds(unknownIDS *map[int]bool) (map[int]string, error) {
	output := map[int]string{}

	ids := make([]int, 0)
	for id := range *unknownIDS {
		ids = append(ids, id)
	}

	// chunk requests underneath the 1,000 item-per-request
	// limit on the eve api
	chunks := util.ChunkSlice(ids, maxItems)
	idm.logger.Debug(fmt.Sprintf("Chunks: %v", chunks))

	for _, chunk := range chunks {
		requestBody, err := json.Marshal(chunk)
		if err != nil {
			return nil, err
		}

		request, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("User-Agent", "hood.noah@icloud.com | github.com/hoodnoah/eve_market/monitor")
		if err != nil {
			return nil, err
		}

		<-idm.rateLimiter.GetChannel()
		response, err := idm.client.Do(request)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download ids %v: %s", chunk, response.Status)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var results []IDResponse
		err = json.Unmarshal(bodyBytes, &results)
		if err != nil {
			return nil, err
		}

		for _, result := range results {
			output[result.ID] = result.Name
		}

	}
	return output, nil
}
