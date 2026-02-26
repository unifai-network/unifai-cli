package errors

import "fmt"

const (
	ExitOK    = 0
	ExitError = 1
	ExitUsage = 2
)

type UsageError struct {
	Message string
}

func (e *UsageError) Error() string {
	return e.Message
}

func NewUsageError(format string, args ...any) error {
	return &UsageError{Message: fmt.Sprintf(format, args...)}
}
