package domain

import "errors"

var (
	ErrValidation                     = errors.New("validation error")
	ErrInternal                       = errors.New("internal error")
	ErrReviewAlreadyExists            = errors.New("review already exists for this course")
	ErrReviewNotFound                 = errors.New("review not found")
	ErrInvalidRating                  = errors.New("rating must be between 1 and 5")
	ErrCourseNotFound                 = errors.New("course not found")
	ErrCourseProgressNotFound         = errors.New("course progress not found")
	ErrCourseSubscriptionNotFound     = errors.New("course subscription not found")
	ErrCertificateNotFound            = errors.New("certificate not found")
	ErrCourseNotCompleted             = errors.New("course is not completed")
	ErrCertificateUnavailable         = errors.New("certificate is not available for this course")
	ErrLessonNotFound                 = errors.New("lesson not found")
	ErrPracticeNotFound               = errors.New("practice not found")
	ErrPracticeSubmissionNotFound     = errors.New("practice submission not found")
	ErrPracticeSubmissionExists       = errors.New("active practice submission already exists")
	ErrPracticeAlreadyCompleted       = errors.New("practice already completed")
	ErrInvalidPracticeReviewStatus    = errors.New("invalid practice review status")
	ErrPracticeManualReviewNotAllowed = errors.New("manual review is not enabled for this practice")
	ErrPracticeAutoSubmitNotAllowed   = errors.New("auto submit is not enabled for this practice")
	ErrQuizNotFound                   = errors.New("quiz not found")
	ErrInactiveUser                   = errors.New("user inactive")
	ErrInvalidToken                   = errors.New("token inactive")
	ErrForbidden                      = errors.New("forbidden")
)
