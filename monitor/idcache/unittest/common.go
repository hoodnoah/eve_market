package idcache_test

import (
	"errors"
	"fmt"

	"github.com/hoodnoah/eve_market/monitor/idcache"
	"github.com/hoodnoah/eve_market/monitor/logger"
)

type TestSetup struct {
	Cache      idcache.IIDcache
	TearDown   func()
	Client     idcache.IRequestClient
	ClientArgs [][]int
}

type FakeClient struct {
	Args          [][]int
	InvalidIDs    []int
	ResponseMap   map[int]string
	Fail          bool
	FailureReason idcache.FailureType
}

// simple linear search
func contains[T comparable](items []T, item T) bool {
	for _, listItem := range items {
		if item == listItem {
			return true
		}
	}

	return false
}

func (fc *FakeClient) FetchManyIDs(ids []int) ([]idcache.ESITypeIDResponse, *idcache.APIRequestError) {
	fc.Args = append(fc.Args, ids)

	if fc.Fail {
		return nil, &idcache.APIRequestError{
			ErrorType: idcache.RequestFailedError,
			Err:       errors.New("request failed for unespecified reason other than 404"),
		}
	}

	results := make([]idcache.ESITypeIDResponse, 0, len(ids))
	for _, id := range ids {
		if fc.ResponseMap[id] != "" {
			record := idcache.ESITypeIDResponse{
				Category: "",
				ID:       id,
				Name:     fc.ResponseMap[id],
			}
			results = append(results, record)
		} else {
			return nil, &idcache.APIRequestError{
				ErrorType: idcache.InvalidIDRequestError,
				Err:       fmt.Errorf("id #%d not valid", id),
			}
		}
	}

	return results, nil
}

func (fc *FakeClient) FetchID(id int) (*idcache.ESITypeIDResponse, *idcache.APIRequestError) {
	fc.Args = append(fc.Args, []int{id})

	if fc.Fail {
		return nil, &idcache.APIRequestError{
			ErrorType: fc.FailureReason,
			Err:       errors.New("fail"),
		}
	}

	if contains(fc.InvalidIDs, id) {
		return nil, &idcache.APIRequestError{
			ErrorType: idcache.InvalidIDRequestError,
			Err:       errors.New(""),
		}
	}

	if fc.ResponseMap[id] == "" {
		err := fmt.Errorf("expected responsemap %v to contain value %d, but it didn't", fc.ResponseMap, id)

		return nil, &idcache.APIRequestError{
			ErrorType: idcache.InvalidRequestError,
			Err:       err,
		}
	}

	return &idcache.ESITypeIDResponse{
		Category: "any",
		ID:       id,
		Name:     fc.ResponseMap[id],
	}, nil
}

func (fc *FakeClient) SetKnownIDs(_ map[int]string) {}

func (fc *FakeClient) Reset() {
	fc.Args = make([][]int, 0)
	fc.InvalidIDs = make([]int, 0)
	fc.ResponseMap = map[int]string{}
}

func (fc *FakeClient) SetInvalidIDS(invalidIDS []int) {
	fc.InvalidIDs = invalidIDS
}

func (fc *FakeClient) SetResponseMap(responseMap map[int]string) {
	fc.ResponseMap = responseMap
}

func (fc *FakeClient) SetFail(fail bool, reason idcache.FailureType) {
	if !fail {
		fc.Fail = false
	} else {
		fc.Fail = true
		fc.FailureReason = reason
	}
}

func newFakeClient() idcache.IRequestClient {
	return &FakeClient{
		Args:        [][]int{},
		ResponseMap: map[int]string{},
		InvalidIDs:  []int{},
	}
}

func setup() *TestSetup {
	consoleLogger := logger.NewConsoleLogger(10)
	consoleLogger.Start()
	client := newFakeClient()

	return &TestSetup{
		Cache:  idcache.NewIDCache(consoleLogger, client),
		Client: client,
	}
}
