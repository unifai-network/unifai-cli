package unifai

import (
	stderrors "errors"
	"fmt"
	"net"
	"net/url"
)

type SearchRequest struct {
	Query          string
	Limit          int
	Offset         int
	IncludeActions []string
}

type InvokeRequest struct {
	Action  string
	Payload any
	Payment any
}

type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("unifai API request failed with status %d", e.StatusCode)
	}
	return fmt.Sprintf("unifai API request failed with status %d: %s", e.StatusCode, e.Body)
}

func IsRetryableError(err error) bool {
	var apiErr *APIError
	if stderrors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500
	}

	var netErr net.Error
	if stderrors.As(err, &netErr) {
		return true
	}

	var urlErr *url.Error
	if stderrors.As(err, &urlErr) {
		return true
	}

	return false
}
