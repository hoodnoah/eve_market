package idcache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/hoodnoah/eve_market/monitor/ratelimiter"
)

type RateLimitedClient struct {
	limiter ratelimiter.IRateLimiter
	client  *http.Client
	mutex   sync.Mutex
}

func NewRateLimitedClient(requestsPerSecond int) IRequestClient {
	rl := ratelimiter.NewTokenBucketRateLimiter(requestsPerSecond)
	rl.Start()

	return &RateLimitedClient{
		client:  &http.Client{},
		limiter: rl,
		mutex:   sync.Mutex{},
	}
}

func (ae *APIRequestError) Error() string {
	return ae.Err.Error()
}

func (rlc *RateLimitedClient) SubmitRequest(request *http.Request) (*http.Response, *APIRequestError) {
	<-rlc.limiter.GetChannel()

	response, err := rlc.client.Do(request)

	if err != nil {
		return nil, &APIRequestError{
			ErrorType: RequestFailedError,
			Err:       err,
		}
	}

	// success
	if response.StatusCode == http.StatusOK {
		return response, nil
	}

	// failure from bad id value
	if response.StatusCode == http.StatusNotFound {
		return nil, &APIRequestError{
			ErrorType: InvalidIDRequestError,
			Err:       nil,
		}
	}

	// failure from any other reason
	return nil, &APIRequestError{
		ErrorType: RequestFailedError,
		Err:       fmt.Errorf("request failed with status %s", response.Status),
	}
}

func (rlc *RateLimitedClient) FetchID(id int) (*ESITypeIDResponse, *APIRequestError) {
	rlc.mutex.Lock()
	defer rlc.mutex.Unlock()

	idList := []int{id}
	body, err := idsToBody(idList)
	if err != nil {
		return nil, &APIRequestError{
			ErrorType: RequestBodyCreationError,
			Err:       err,
		}
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, &APIRequestError{
			ErrorType: RequestCreationError,
			Err:       err,
		}
	}

	res, resError := rlc.SubmitRequest(req)
	if resError != nil {
		switch resError.ErrorType {
		case InvalidIDRequestError:
			unknownNameVal := fmt.Sprintf("unknownID_%d", id)
			return &ESITypeIDResponse{
				Category: "",
				ID:       id,
				Name:     unknownNameVal,
			}, nil
		default:
			return nil, resError
		}
	}

	resBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &APIRequestError{
			ErrorType: ResponseParseError,
			Err:       err,
		}
	}

	// unmarshal results
	var results []ESITypeIDResponse
	err = json.Unmarshal(resBodyBytes, &results)
	if err != nil {
		return nil, &APIRequestError{
			ErrorType: ResponseParseError,
			Err:       err,
		}
	}
	if len(results) != 1 {
		err = fmt.Errorf("expected only one response from the EVE API, received %d", len(results))
		return nil, &APIRequestError{
			ErrorType: ResponseParseError,
			Err:       err,
		}
	}

	return &ESITypeIDResponse{
		Category: results[0].Category,
		ID:       results[0].ID,
		Name:     results[0].Name,
	}, nil
}

func (rlc *RateLimitedClient) FetchManyIDs(ids []int) ([]ESITypeIDResponse, *APIRequestError) {
	// make the body
	requestBody, err := idsToBody(ids)
	if err != nil {
		return nil, &APIRequestError{
			ErrorType: RequestBodyCreationError,
			Err:       err,
		}
	}

	request, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return nil, &APIRequestError{
			ErrorType: RequestCreationError,
			Err:       err,
		}
	}

	response, requestErr := rlc.SubmitRequest(request)
	if requestErr != nil {
		return nil, requestErr
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, &APIRequestError{
			ErrorType: InvalidIDRequestError,
			Err:       fmt.Errorf("failed to fetch ids: %s", response.Status),
		}
	}

	if response.StatusCode != http.StatusOK {
		return nil, &APIRequestError{
			ErrorType: RequestFailedError,
			Err:       fmt.Errorf("failed to fetch ids: %s", response.Status),
		}
	}

	resBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &APIRequestError{
			ErrorType: ResponseParseError,
			Err:       err,
		}
	}

	var esiResults []ESITypeIDResponse
	err = json.Unmarshal(resBodyBytes, &esiResults)

	if err != nil {
		return nil, &APIRequestError{
			ErrorType: ResponseParseError,
			Err:       err,
		}
	}

	return esiResults, nil
}

func idsToBody(ids []int) (io.Reader, error) {
	jsonValue, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonValue), nil
}
