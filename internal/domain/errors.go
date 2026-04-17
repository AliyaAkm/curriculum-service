package domain

import "errors"

var (
	ErrValidation          = errors.New("validation error")
	ErrInternal            = errors.New("internal error")
	ErrReviewAlreadyExists = errors.New("review already exists for this course")
	ErrReviewNotFound      = errors.New("review not found")
	ErrInvalidRating       = errors.New("rating must be between 1 and 5")
	ErrCourseNotFound      = errors.New("course not found")
)
