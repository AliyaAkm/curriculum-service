package course

import (
	"context"
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/course"
	"curriculum-service/internal/domain/review"
	"curriculum-service/internal/domain/reviewlog"
	dtocourse "curriculum-service/internal/http/dto/course"
	"curriculum-service/internal/http/dto/resubmit"
	"curriculum-service/internal/http/dto/reviewcourse"
	"curriculum-service/internal/http/dto/reviewresult"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

func (u *UseCase) GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error) {
	resp, err := u.repo.GetAllCourses(ctx, query)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(resp); i++ {
		reviews, err := u.reviewRepo.GetAllReviewsByCourseID(ctx, resp[i].ID)
		if err != nil {
			return nil, err
		}

		resp[i].Rating, resp[i].RatingCount = calculateRatingStats(reviews)
	}

	return resp, nil
}

func (u *UseCase) GetPendingCheckCourses(ctx context.Context) ([]course.Course, error) {
	return u.repo.GetPendingCheckCourses(ctx)
}

func (u *UseCase) ReviewCourse(
	ctx context.Context,
	courseID, adminID uuid.UUID,
	req reviewcourse.ReviewCourseRequest,
) (*reviewcourse.ReviewCourseResponse, error) {
	comment := strings.TrimSpace(req.Comment)


	log := &reviewlog.CourseReviewLog{
		ID:        uuid.New(),
		CourseID:  courseID,
		AdminID:   adminID,
		IsApproved: req.IsApproved,
		Comment:   comment,
		CreatedAt: time.Now(),
	}

	if err := u.repo.ReviewCourse(ctx, log, req.IsApproved); err != nil {
		return nil, err
	}

	return &reviewcourse.ReviewCourseResponse{
		CourseID:   courseID,
		IsChecked:  true,
		IsApproved: req.IsApproved,
		Comment:    log.Comment,
		ReviewedBy: adminID,
		ReviewedAt: log.CreatedAt,
	}, nil
}

func reviewStatus(isChecked, isApproved bool) string {
	switch {
	case !isChecked:
		return "pending_review"
	case isApproved:
		return "approved"
	default:
		return "rejected"
	}
}

func (u *UseCase) GetCourseReview(
	ctx context.Context,
	courseID, userID uuid.UUID,
) (*reviewresult.CourseReviewResponse, error) {
	c, err := u.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // adapt to GetCourseByID's miss behavior
			return nil, domain.ErrCourseNotFound
		}
		return nil, err
	}

	// owner OR (admin allowed to view all reviews) — otherwise 403
	if c.AuthorID != userID {
		return nil, domain.ErrNotCourseOwner
	}

	resp := &reviewresult.CourseReviewResponse{
		CourseID:     c.ID,
		Title:        c.Title,
		IsChecked:    c.IsChecked,
		IsApproved:   c.IsApproved,
		ReviewStatus: reviewStatus(c.IsChecked, c.IsApproved),
	}

	// latest admin comment, only if the review-log history exists
	log, err := u.repo.GetLatestReviewLog(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if log != nil {
		if log.Comment != "" {
			resp.Comment = &log.Comment
		}
		resp.ReviewedBy = &log.AdminID
		resp.ReviewedAt = &log.CreatedAt
	}

	return resp, nil
}

func (u *UseCase) ResubmitCourseForReview(
	ctx context.Context,
	courseID, userID uuid.UUID,
) (*resubmit.ResubmitReviewResponse, error) {
	c, err := u.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		// adapt this to whatever GetCourseByID returns on miss
		// (gorm.ErrRecordNotFound or your existing ErrCourseNotFound)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCourseNotFound
		}
		return nil, err
	}

	// ownership check using the existing author field
	if c.AuthorID != userID {
		return nil, domain.ErrNotCourseOwner
	}

	if err := u.repo.SetCourseUnchecked(ctx, courseID); err != nil {
		return nil, err
	}

	return &resubmit.ResubmitReviewResponse{
		CourseID:   courseID,
		IsChecked:  false,
		IsApproved: false,
		Message:    "Course has been sent back for admin review",
	}, nil
}

func (u *UseCase) CreateCourse(ctx context.Context, value *course.Course) (*course.Course, error) {
	value.ID = uuid.New()
	value.Rating = 0
	value.RatingCount = 0

	id, err := u.repo.CreateCourse(ctx, value)
	if err != nil {
		return nil, err
	}

	return u.repo.GetCourseByID(ctx, id)
}

func (u *UseCase) CreateSubscription(ctx context.Context, value *course.Subscription) (*course.Subscription, error) {
	value.ID = uuid.New()
	err := u.repo.CreateSubscription(ctx, value)
	if err != nil {
		return nil, err
	}
	subscription, err := u.repo.GetSubscriptionByID(ctx, value.ID)
	if err != nil {
		return nil, err
	}

	if u.notification != nil {
		data := map[string]any{
			"courseId":       value.CourseID.String(),
			"subscriptionId": value.ID.String(),
		}
		if courseEntity, err := u.repo.GetCourseByID(ctx, value.CourseID); err == nil && courseEntity != nil {
			data["courseTitle"] = courseEntity.Title
		}
		_ = u.notification.SendEvent(ctx, value.UserID, "course_enrolled", data)
	}

	return subscription, nil
}

func (u *UseCase) GetCourseForUser(ctx context.Context, userID uuid.UUID, courseID uuid.UUID, hasFullAccess bool) (*course.CourseForUser, error) {
	courseEntity, err := u.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	subscription, err := u.repo.HasSubscription(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}

	modules, err := u.moduleRepo.GetModuleByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return &course.CourseForUser{
		Course:          courseEntity,
		Modules:         modules,
		HasSubscription: subscription,
	}, nil
}

func (u *UseCase) GetCourseByID(ctx context.Context, id uuid.UUID) (*course.Course, error) {
	courseValue, err := u.repo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	reviews, err := u.reviewRepo.GetAllReviewsByCourseID(ctx, courseValue.ID)
	if err != nil {
		return nil, err
	}

	courseValue.Rating, courseValue.RatingCount = calculateRatingStats(reviews)

	return courseValue, nil
}

func (u *UseCase) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	err := u.repo.DeleteCoursePrice(ctx, id)
	if err != nil {
		return err
	}
	return u.repo.DeleteCourse(ctx, id)
}
func (u *UseCase) UpdateCourse(ctx context.Context, id uuid.UUID, value *course.Course) (*course.Course, error) {
	err := u.repo.UpdateCourse(ctx, id, value)
	if err != nil {
		return nil, err
	}

	courseValue, err := u.repo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	reviews, err := u.reviewRepo.GetAllReviewsByCourseID(ctx, courseValue.ID)
	if err != nil {
		return nil, err
	}

	courseValue.Rating, courseValue.RatingCount = calculateRatingStats(reviews)

	return courseValue, nil
}

func calculateRatingStats(reviews []review.CourseReview) (float64, int) {
	if len(reviews) == 0 {
		return 0, 0
	}

	var sum float64
	for i := 0; i < len(reviews); i++ {
		sum += float64(reviews[i].Rating)
	}

	avg := sum / float64(len(reviews))
	return math.Round(avg*10) / 10, len(reviews)
}
