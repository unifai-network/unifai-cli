package retry

import (
	"context"
	"time"
)

func Do[T any](ctx context.Context, maxRetries int, initialBackoff time.Duration, shouldRetry func(error) bool, fn func() (T, error)) (T, error) {
	var zero T
	if maxRetries < 0 {
		maxRetries = 0
	}
	if initialBackoff <= 0 {
		initialBackoff = time.Second
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}

		if attempt == maxRetries || (shouldRetry != nil && !shouldRetry(err)) {
			return zero, err
		}

		backoff := initialBackoff << attempt
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(backoff):
		}
	}

	// unreachable: loop always returns on success, final attempt, or context cancellation
	return zero, nil
}
