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

// ActiveTasksLimitError
type ActiveTasksLimitError struct {
	MaxAllowed int
}

func (e ActiveTasksLimitError) Error() string {
	return fmt.Sprintf("active tasks limit exceeded: %d allowed",
		e.MaxAllowed)
}

func IsActiveTasksLimitError(err error) bool {
	_, ok := err.(ActiveTasksLimitError)
	return ok
}

// NotFoundError
type NotFoundError struct {
	Entity string
	ID     string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %s not found",
		e.Entity, e.ID)
}

func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

// ArchiveReadyError
type ArchiveReadyError struct {
	ID string
}

func (e ArchiveReadyError) Error() string {
	return fmt.Sprintf("The task with id %s is already archived",
		e.ID)
}

func IsArchiveReadyError(err error) bool {
	_, ok := err.(ArchiveReadyError)
	return ok
}
