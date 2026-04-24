package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	"curriculum-service/internal/domain/review"
	dtocourse "curriculum-service/internal/http/dto/course"
	"github.com/google/uuid"
	"math"
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
	return u.repo.GetSubscriptionByID(ctx, value.ID)
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
