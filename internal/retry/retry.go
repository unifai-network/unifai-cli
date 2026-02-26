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

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		lastErr = err

		if attempt == maxRetries || (shouldRetry != nil && !shouldRetry(err)) {
			return zero, err
		}

		backoff := initialBackoff << attempt
		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, ctx.Err()
		case <-timer.C:
		}
	}

	return zero, lastErr
}
