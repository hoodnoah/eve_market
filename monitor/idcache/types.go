package idcache

import (
	"sync"

	"github.com/hoodnoah/eve_market/monitor/logger"
)

type IDType int

const (
	RegionID IDType = iota
	TypeID
)

// interface
type IIDcache interface {
	SetKnownIDs(ids map[int]string)
	Label(id int) (string, error)
	LabelMany(ids []int) (map[int]string, error)
}

type IRequestClient interface {
	FetchID(id int) (*ESITypeIDResponse, *APIRequestError)
	FetchManyIDs(id []int) ([]ESITypeIDResponse, *APIRequestError)
}

type IDCache struct {
	logger logger.ILogger
	client IRequestClient
	ids    map[int]string
	mutex  sync.Mutex
}

type ESITypeIDResponse struct {
	Category string `json:"category"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
}

type FailureType int

const (
	RequestBodyCreationError FailureType = iota
	RequestCreationError
	RequestFailedError
	InvalidRequestError
	InvalidIDRequestError
	ResponseParseError
)

type APIRequestError struct {
	ErrorType FailureType
	Err       error
}
