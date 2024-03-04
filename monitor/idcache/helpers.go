package idcache

import (
	"fmt"
	"net/http"
)

var retryableErrors = []int{
	500,
	503,
	504,
}

func isIn[T comparable](list []T, value T) bool {
	for _, item := range list {
		if value == item {
			return true
		}
	}
	return false
}

func isRetryableError(response *http.Response) bool {
	return isIn(retryableErrors, response.StatusCode)
}

// handles the response from the eve ESI API
func handleAPIError(err error, response *http.Response) *APIRequestError {
	// request failed outright
	if err != nil {
		return &APIRequestError{
			ErrorType: RequestFailedError,
			Err:       err,
		}
	}

	if response == nil {
		return &APIRequestError{
			ErrorType: RequestFailedError,
			Err:       fmt.Errorf("request failed, nil response returned"),
		}
	}

	// request was successful
	if response.StatusCode == http.StatusOK {
		return nil
	}

	// request failed due to invalid if
	if response.StatusCode == http.StatusNotFound {
		return &APIRequestError{
			ErrorType: InvalidIDRequestError,
			Err:       err,
		}
	}

	// request failed for any other reason
	return &APIRequestError{
		ErrorType: RequestFailedError,
		Err:       fmt.Errorf("request failed with status %s", response.Status),
	}
}
