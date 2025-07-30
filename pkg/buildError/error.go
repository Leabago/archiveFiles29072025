package builderror

import (
	"errors"
	"fmt"
)

// добавить параметр в ошибку
func WithParam(err error, key string, value interface{}) error {
	return fmt.Errorf("%s=%v: Error: %w", key, value, err)
}

func ErrUrl(message, url string) error {
	fail := errors.New(message)
	fail = WithParam(fail, "url", url)
	return fail
}

// LinkLimitError
type LinkLimitError struct {
	Requested  int
	MaxAllowed int
}

func (e LinkLimitError) Error() string {
	return fmt.Sprintf("link limit exceeded: %d allowed, but got %d",
		e.MaxAllowed, e.Requested)
}
func IsLinkLimitError(err error) bool {
	_, ok := err.(LinkLimitError)
	return ok
}
