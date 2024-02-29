package idcache

import (
	"net/http"
	"sync"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/ratelimiter"
)

type IIDCache interface {
	SetKnownRegionIDs(regionIDs *RegionIDInput)
	SetKnownTypeIDs(typeIDs *TypeIDInput)
	Label(ids *UnknownIDs) (*KnownIDs, error)
}

type IDType int

const (
	RegionID IDType = iota
	TypeID
)

type IDValue struct {
	ID    int
	Value string
}

type UnknownIDs struct {
	Type IDType
	IDS  map[int]bool
}

type KnownIDs struct {
	Type IDType
	IDS  map[int]string
}

type RegionIDInput KnownIDs
type TypeIDInput KnownIDs

type IDCache struct {
	logger      logger.ILogger
	regionIDS   map[int]string
	typeIDS     map[int]string
	mutex       sync.Mutex
	rateLimiter ratelimiter.IRateLimiter
	client      http.Client
}

type IDRequest struct {
	IDs []int
}

type IDResponse struct {
	Category string `json:"category"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}
