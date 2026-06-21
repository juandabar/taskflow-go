package apperror

import "fmt"

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' does not exist", e.Resource, e.ID)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{Resource: resource, ID: id}
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{Message: message}
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}
