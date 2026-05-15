package quiz

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"curriculum-service/internal/domain"
	quizdomain "curriculum-service/internal/domain/quiz"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type quizRow struct {
	ID                 string
	LessonID           string
	Position           int
	CorrectAnswerIndex int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type quizTextRow struct {
	QuizID      string
	Code        string
	Question    string
	Explanation string
}

type optionRow struct {
	ID       string
	QuizID   string
	Position int
}

type optionTextRow struct {
	OptionID string
	Code     string
	Text     string
}

type localeRow struct {
	ID   string
	Code string
}

func (r *Repo) CreateQuiz(ctx context.Context, value *quizdomain.Quiz) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		locales, err := getLocaleIDs(ctx, tx)
		if err != nil {
			return err
		}

		if err := tx.Exec(`
			INSERT INTO lesson_quizzes (
				id, lesson_id, position, correct_answer_index
			) VALUES (?, ?, ?, ?)
		`, value.ID, value.LessonID, value.Position, value.CorrectAnswerIndex).Error; err != nil {
			return err
		}

		if err := createQuizTexts(ctx, tx, locales, value.ID, value.Question, value.Explanation); err != nil {
			return err
		}

		return createOptions(ctx, tx, locales, value.ID, value.Options)
	})
}

func (r *Repo) GetQuizByID(ctx context.Context, id uuid.UUID) (*quizdomain.Quiz, error) {
	quizzes, err := r.getQuizzes(ctx, "q.id = ?", id)
	if err != nil {
		return nil, err
	}
	if len(quizzes) == 0 {
		return nil, domain.ErrQuizNotFound
	}

	return &quizzes[0], nil
}

func (r *Repo) GetQuizzesByLessonID(ctx context.Context, lessonID uuid.UUID) ([]quizdomain.Quiz, error) {
	return r.getQuizzes(ctx, "q.lesson_id = ?", lessonID)
}

func (r *Repo) UpdateQuiz(ctx context.Context, id uuid.UUID, value *quizdomain.Quiz) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
			UPDATE lesson_quizzes
			SET position = ?,
			    correct_answer_index = ?,
			    updated_at = NOW()
			WHERE id = ?
		`, value.Position, value.CorrectAnswerIndex, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrQuizNotFound
		}

		locales, err := getLocaleIDs(ctx, tx)
		if err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Exec(`DELETE FROM lesson_quiz_texts WHERE quiz_id = ?`, id).Error; err != nil {
			return err
		}
		if err := tx.WithContext(ctx).Exec(`DELETE FROM lesson_quiz_options WHERE quiz_id = ?`, id).Error; err != nil {
			return err
		}
		if err := createQuizTexts(ctx, tx, locales, id, value.Question, value.Explanation); err != nil {
			return err
		}

		return createOptions(ctx, tx, locales, id, value.Options)
	})
}

func (r *Repo) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Exec(`DELETE FROM lesson_quizzes WHERE id = ?`, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrQuizNotFound
	}

	return nil
}

func (r *Repo) SaveQuizAttempt(ctx context.Context, userID uuid.UUID, quizID uuid.UUID, selectedAnswerIndex int, isCorrect bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(`
			UPDATE lesson_quiz_attempts
			SET selected_answer_index = ?,
			    is_correct = ?
			WHERE user_id = ?
			  AND quiz_id = ?
		`, selectedAnswerIndex, isCorrect, userID, quizID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			return nil
		}

		return tx.Exec(`
			INSERT INTO lesson_quiz_attempts (
				id, quiz_id, user_id, selected_answer_index, is_correct
			) VALUES (?, ?, ?, ?, ?)
		`, uuid.New(), quizID, userID, selectedAnswerIndex, isCorrect).Error
	})
}

func (r *Repo) GetLessonAccessInfo(ctx context.Context, lessonID uuid.UUID) (uuid.UUID, uuid.UUID, int, error) {
	var row struct {
		CourseID       string `gorm:"column:course_id"`
		ModuleID       string `gorm:"column:module_id"`
		ModulePosition int    `gorm:"column:module_position"`
	}

	err := r.db.WithContext(ctx).
		Table("course_lessons cl").
		Select("cm.course_id AS course_id, cl.module_id AS module_id, COALESCE(cm.position, 0) AS module_position").
		Joins("INNER JOIN course_modules cm ON cm.id = cl.module_id").
		Where("cl.id = ?", lessonID).
		Scan(&row).Error
	if err != nil {
		return uuid.Nil, uuid.Nil, 0, err
	}

	courseID, err := uuid.Parse(row.CourseID)
	if err != nil {
		return uuid.Nil, uuid.Nil, 0, err
	}
	moduleID, err := uuid.Parse(row.ModuleID)
	if err != nil {
		return uuid.Nil, uuid.Nil, 0, err
	}

	return courseID, moduleID, row.ModulePosition, nil
}

func (r *Repo) HasSubscription(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("course_subscription").
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) IsModuleInFreePreview(ctx context.Context, courseID uuid.UUID, moduleID uuid.UUID, limit int) (bool, error) {
	var allowed bool
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT EXISTS (
				SELECT 1
				FROM course_modules
				WHERE course_id = ?
				  AND id = ?
				  AND COALESCE(position, 0) IN (
					SELECT position
					FROM course_modules
					WHERE course_id = ?
					  AND position IS NOT NULL
					GROUP BY position
					ORDER BY position ASC
					LIMIT ?
				  )
			)
		`, courseID, moduleID, courseID, limit).
		Scan(&allowed).Error
	if err != nil {
		return false, err
	}
	return allowed, nil
}

func (r *Repo) getQuizzes(ctx context.Context, condition string, args ...any) ([]quizdomain.Quiz, error) {
	var rows []quizRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			q.id::text,
			q.lesson_id::text,
			q.position,
			q.correct_answer_index,
			q.created_at,
			q.updated_at
		FROM lesson_quizzes q
		WHERE `+condition+`
		ORDER BY q.position ASC, q.created_at ASC
	`, args...).Scan(&rows).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrQuizNotFound
		}
		return nil, err
	}

	quizzes := make([]quizdomain.Quiz, len(rows))
	quizIDs := make([]uuid.UUID, len(rows))
	for i := range rows {
		quizEntity, err := convertQuizRow(rows[i])
		if err != nil {
			return nil, err
		}
		quizzes[i] = quizEntity
		quizIDs[i] = quizEntity.ID
	}
	if len(quizzes) == 0 {
		return quizzes, nil
	}

	if err := r.loadTexts(ctx, quizzes, quizIDs); err != nil {
		return nil, err
	}
	if err := r.loadOptions(ctx, quizzes, quizIDs); err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (r *Repo) loadTexts(ctx context.Context, quizzes []quizdomain.Quiz, quizIDs []uuid.UUID) error {
	var rows []quizTextRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			qt.quiz_id::text,
			cl.code,
			qt.question,
			qt.explanation
		FROM lesson_quiz_texts qt
		INNER JOIN course_locales cl ON cl.id = qt.locale_id
		WHERE qt.quiz_id IN ?
	`, quizIDs).Scan(&rows).Error
	if err != nil {
		return err
	}

	indexByID := make(map[string]int, len(quizzes))
	for i := range quizzes {
		indexByID[quizzes[i].ID.String()] = i
	}
	for _, row := range rows {
		i, ok := indexByID[row.QuizID]
		if !ok {
			continue
		}
		setLocaleValue(&quizzes[i].Question, row.Code, row.Question)
		setLocaleValue(&quizzes[i].Explanation, row.Code, row.Explanation)
	}

	return nil
}

func (r *Repo) loadOptions(ctx context.Context, quizzes []quizdomain.Quiz, quizIDs []uuid.UUID) error {
	var optionRows []optionRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT id::text, quiz_id::text, position
		FROM lesson_quiz_options
		WHERE quiz_id IN ?
		ORDER BY quiz_id ASC, position ASC
	`, quizIDs).Scan(&optionRows).Error
	if err != nil {
		return err
	}
	if len(optionRows) == 0 {
		return nil
	}

	optionIDs := make([]uuid.UUID, 0, len(optionRows))
	quizIndexByID := make(map[string]int, len(quizzes))
	optionIndexByID := make(map[string]struct {
		quizIndex   int
		optionIndex int
	}, len(optionRows))
	for i := range quizzes {
		quizIndexByID[quizzes[i].ID.String()] = i
	}
	for _, row := range optionRows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return err
		}
		quizIndex, ok := quizIndexByID[row.QuizID]
		if !ok {
			continue
		}
		quizzes[quizIndex].Options = append(quizzes[quizIndex].Options, quizdomain.Option{
			ID:       id,
			Position: row.Position,
		})
		optionIDs = append(optionIDs, id)
		optionIndexByID[id.String()] = struct {
			quizIndex   int
			optionIndex int
		}{
			quizIndex:   quizIndex,
			optionIndex: len(quizzes[quizIndex].Options) - 1,
		}
	}

	var textRows []optionTextRow
	err = r.db.WithContext(ctx).Raw(`
		SELECT
			ot.option_id::text,
			cl.code,
			ot.text
		FROM lesson_quiz_option_texts ot
		INNER JOIN course_locales cl ON cl.id = ot.locale_id
		WHERE ot.option_id IN ?
	`, optionIDs).Scan(&textRows).Error
	if err != nil {
		return err
	}

	for _, row := range textRows {
		location, ok := optionIndexByID[row.OptionID]
		if !ok {
			continue
		}
		setLocaleValue(&quizzes[location.quizIndex].Options[location.optionIndex].Text, row.Code, row.Text)
	}

	return nil
}

func getLocaleIDs(ctx context.Context, tx *gorm.DB) (map[string]uuid.UUID, error) {
	var rows []localeRow
	if err := tx.WithContext(ctx).Raw(`SELECT id::text, code FROM course_locales`).Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make(map[string]uuid.UUID, len(rows))
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return nil, err
		}
		result[row.Code] = id
	}

	for _, code := range []string{"en", "ru", "kk"} {
		if result[code] == uuid.Nil {
			return nil, fmt.Errorf("%w: missing course locale %s", domain.ErrValidation, code)
		}
	}

	return result, nil
}

func createQuizTexts(ctx context.Context, tx *gorm.DB, locales map[string]uuid.UUID, quizID uuid.UUID, question quizdomain.Locale, explanation quizdomain.Locale) error {
	for _, item := range []struct {
		code        string
		question    string
		explanation string
	}{
		{code: "en", question: question.EN, explanation: explanation.EN},
		{code: "ru", question: question.RU, explanation: explanation.RU},
		{code: "kk", question: question.KK, explanation: explanation.KK},
	} {
		if err := tx.WithContext(ctx).Exec(`
			INSERT INTO lesson_quiz_texts (
				id, quiz_id, locale_id, question, explanation
			) VALUES (?, ?, ?, ?, ?)
		`, uuid.New(), quizID, locales[item.code], item.question, item.explanation).Error; err != nil {
			return err
		}
	}

	return nil
}

func createOptions(ctx context.Context, tx *gorm.DB, locales map[string]uuid.UUID, quizID uuid.UUID, options []quizdomain.Option) error {
	for _, option := range options {
		if err := tx.WithContext(ctx).Exec(`
			INSERT INTO lesson_quiz_options (
				id, quiz_id, position
			) VALUES (?, ?, ?)
		`, option.ID, quizID, option.Position).Error; err != nil {
			return err
		}

		for _, item := range []struct {
			code string
			text string
		}{
			{code: "en", text: option.Text.EN},
			{code: "ru", text: option.Text.RU},
			{code: "kk", text: option.Text.KK},
		} {
			if err := tx.WithContext(ctx).Exec(`
				INSERT INTO lesson_quiz_option_texts (
					id, option_id, locale_id, text
				) VALUES (?, ?, ?, ?)
			`, uuid.New(), option.ID, locales[item.code], item.text).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func convertQuizRow(row quizRow) (quizdomain.Quiz, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return quizdomain.Quiz{}, err
	}
	lessonID, err := uuid.Parse(row.LessonID)
	if err != nil {
		return quizdomain.Quiz{}, err
	}

	return quizdomain.Quiz{
		ID:                 id,
		LessonID:           lessonID,
		Position:           row.Position,
		CorrectAnswerIndex: row.CorrectAnswerIndex,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}, nil
}

func setLocaleValue(dst *quizdomain.Locale, localeCode string, value string) {
	switch localeCode {
	case "en":
		dst.EN = value
	case "ru":
		dst.RU = value
	case "kk":
		dst.KK = value
	}
}
