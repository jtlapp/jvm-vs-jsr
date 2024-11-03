package command

import "fmt"

type UsageError struct {
	message string
}

func (e *UsageError) Error() string {
	return e.message
}

func NewUsageError(format string, args ...interface{}) *UsageError {
	return &UsageError{message: fmt.Sprintf(format, args...)}
}

func IsUsageError(err error) bool {
	_, ok := err.(*UsageError)
	return ok
}
