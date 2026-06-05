package quiz

import (
	"context"
	"curriculum-service/internal/domain"
	achievementdomain "curriculum-service/internal/domain/achievement"
	quizdomain "curriculum-service/internal/domain/quiz"

	"github.com/google/uuid"
)

const freePreviewModulesLimit = 2

func (u *UseCase) CreateQuiz(ctx context.Context, value *quizdomain.Quiz) (*quizdomain.Quiz, error) {
	if err := prepareQuiz(value); err != nil {
		return nil, err
	}

	value.ID = uuid.New()
	for i := range value.Options {
		value.Options[i].ID = uuid.New()
		value.Options[i].Position = i + 1
	}

	if err := u.repo.CreateQuiz(ctx, value); err != nil {
		return nil, err
	}

	return u.repo.GetQuizByID(ctx, value.ID)
}

func (u *UseCase) GetQuizzesByLessonIDForUser(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) ([]quizdomain.Quiz, error) {
	if err := u.requireLessonAccess(ctx, userID, lessonID, hasFullAccess); err != nil {
		return nil, err
	}

	return u.repo.GetQuizzesByLessonID(ctx, lessonID)
}

func (u *UseCase) GetQuizByIDForUser(ctx context.Context, userID uuid.UUID, id uuid.UUID, hasFullAccess bool) (*quizdomain.Quiz, error) {
	quizEntity, err := u.repo.GetQuizByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := u.requireLessonAccess(ctx, userID, quizEntity.LessonID, hasFullAccess); err != nil {
		return nil, err
	}

	return quizEntity, nil
}

func (u *UseCase) SubmitAnswer(ctx context.Context, userID uuid.UUID, quizID uuid.UUID, selectedAnswerIndex int, hasFullAccess bool) (*quizdomain.AnswerResult, error) {
	quizEntity, err := u.repo.GetQuizByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	if err := u.requireLessonAccess(ctx, userID, quizEntity.LessonID, hasFullAccess); err != nil {
		return nil, err
	}
	if selectedAnswerIndex < 0 || selectedAnswerIndex >= len(quizEntity.Options) {
		return nil, domain.ErrValidation
	}

	isCorrect := selectedAnswerIndex == quizEntity.CorrectAnswerIndex
	if err := u.repo.SaveQuizAttempt(ctx, userID, quizID, selectedAnswerIndex, isCorrect); err != nil {
		return nil, err
	}

	var correctOptionID uuid.UUID
	if quizEntity.CorrectAnswerIndex >= 0 && quizEntity.CorrectAnswerIndex < len(quizEntity.Options) {
		correctOptionID = quizEntity.Options[quizEntity.CorrectAnswerIndex].ID
	}

	if u.notification != nil {
		score := "0%"
		if isCorrect {
			score = "100%"
		}
		_ = u.notification.SendEvent(ctx, userID, "assessment_completed", map[string]any{
			"score":    score,
			"quizId":   quizEntity.ID.String(),
			"lessonId": quizEntity.LessonID.String(),
		})
		if isCorrect {
			u.notifyUnlockedAchievements(ctx, userID)
		}
	}

	return &quizdomain.AnswerResult{
		QuizID:              quizEntity.ID,
		SelectedAnswerIndex: selectedAnswerIndex,
		IsCorrect:           isCorrect,
		CorrectAnswerIndex:  quizEntity.CorrectAnswerIndex,
		CorrectOptionID:     correctOptionID,
		Explanation:         quizEntity.Explanation,
	}, nil
}

func (u *UseCase) UpdateQuiz(ctx context.Context, id uuid.UUID, value *quizdomain.Quiz) (*quizdomain.Quiz, error) {
	if err := prepareQuizContent(value); err != nil {
		return nil, err
	}
	for i := range value.Options {
		value.Options[i].ID = uuid.New()
		value.Options[i].Position = i + 1
	}

	if err := u.repo.UpdateQuiz(ctx, id, value); err != nil {
		return nil, err
	}

	return u.repo.GetQuizByID(ctx, id)
}

func (u *UseCase) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteQuiz(ctx, id)
}

func (u *UseCase) requireLessonAccess(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) error {
	if hasFullAccess {
		return nil
	}

	courseID, moduleID, _, err := u.repo.GetLessonAccessInfo(ctx, lessonID)
	if err != nil {
		return err
	}

	hasSubscription, err := u.repo.HasSubscription(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if hasSubscription {
		return nil
	}

	allowed, err := u.repo.IsModuleInFreePreview(ctx, courseID, moduleID, freePreviewModulesLimit)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.ErrForbidden
	}

	return nil
}

func (u *UseCase) notifyUnlockedAchievements(ctx context.Context, userID uuid.UUID) {
	if u.notification == nil || u.achievements == nil {
		return
	}

	items, err := u.achievements.SyncUnlockedAchievements(ctx, userID)
	if err != nil {
		return
	}
	for _, item := range items {
		_ = u.notification.SendEvent(ctx, userID, "achievement_unlocked", map[string]any{
			"achievementId":    item.ID.String(),
			"achievementCode":  item.Code,
			"achievementTitle": achievementTitle(item),
		})
	}
}

func achievementTitle(value achievementdomain.Achievement) string {
	switch {
	case value.Title.RU != "":
		return value.Title.RU
	case value.Title.EN != "":
		return value.Title.EN
	case value.Title.KK != "":
		return value.Title.KK
	default:
		return "новое достижение"
	}
}

func prepareQuiz(value *quizdomain.Quiz) error {
	if value == nil || value.LessonID == uuid.Nil {
		return domain.ErrValidation
	}
	return prepareQuizContent(value)
}

func prepareQuizContent(value *quizdomain.Quiz) error {
	if value == nil {
		return domain.ErrValidation
	}
	if value.Position <= 0 {
		value.Position = 1
	}
	if len(value.Options) < 2 {
		return domain.ErrValidation
	}
	if value.CorrectAnswerIndex < 0 || value.CorrectAnswerIndex >= len(value.Options) {
		return domain.ErrValidation
	}

	return nil
}
