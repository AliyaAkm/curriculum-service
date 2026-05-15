package practice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"curriculum-service/internal/domain"
	practicedomain "curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
)

type practiceRow struct {
	ID                     string
	LessonID               string
	Position               int
	TitleEN                string
	TitleRU                string
	TitleKK                string
	SummaryEN              string
	SummaryRU              string
	SummaryKK              string
	BriefEN                string
	BriefRU                string
	BriefKK                string
	StarterCode            string
	SuccessCriteriaJSONRaw []byte
	KnowledgeChecksJSONRaw []byte
	PromptSuggestionEN     string
	PromptSuggestionRU     string
	PromptSuggestionKK     string
	XPReward               int
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

func (r *Repo) CreatePractice(ctx context.Context, value *practicedomain.Practice) error {
	successCriteriaJSON, err := marshalLocales(value.SuccessCriteria)
	if err != nil {
		return err
	}
	knowledgeChecksJSON, err := marshalLocales(value.KnowledgeChecks)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Exec(`
		INSERT INTO course_practices (
			id, lesson_id, position,
			title_en, title_ru, title_kk,
			summary_en, summary_ru, summary_kk,
			brief_en, brief_ru, brief_kk,
			starter_code, success_criteria, knowledge_checks,
			prompt_suggestion_en, prompt_suggestion_ru, prompt_suggestion_kk,
			xp_reward
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?::jsonb, ?::jsonb, ?, ?, ?, ?)
	`, value.ID, value.LessonID, value.Position,
		value.Title.EN, value.Title.RU, value.Title.KK,
		value.Summary.EN, value.Summary.RU, value.Summary.KK,
		value.Brief.EN, value.Brief.RU, value.Brief.KK,
		value.StarterCode, successCriteriaJSON, knowledgeChecksJSON,
		value.PromptSuggestion.EN, value.PromptSuggestion.RU, value.PromptSuggestion.KK,
		value.XPReward,
	).Error
}

func (r *Repo) GetPracticeByID(ctx context.Context, id uuid.UUID) (*practicedomain.Practice, error) {
	return r.getPractice(ctx, "id = ?", id)
}

func (r *Repo) GetPracticeByLessonID(ctx context.Context, lessonID uuid.UUID) (*practicedomain.Practice, error) {
	return r.getPractice(ctx, "lesson_id = ?", lessonID)
}

func (r *Repo) UpdatePractice(ctx context.Context, id uuid.UUID, value *practicedomain.Practice) error {
	successCriteriaJSON, err := marshalLocales(value.SuccessCriteria)
	if err != nil {
		return err
	}
	knowledgeChecksJSON, err := marshalLocales(value.KnowledgeChecks)
	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Exec(`
		UPDATE course_practices
		SET lesson_id = ?,
		    position = ?,
		    title_en = ?,
		    title_ru = ?,
		    title_kk = ?,
		    summary_en = ?,
		    summary_ru = ?,
		    summary_kk = ?,
		    brief_en = ?,
		    brief_ru = ?,
		    brief_kk = ?,
		    starter_code = ?,
		    success_criteria = ?::jsonb,
		    knowledge_checks = ?::jsonb,
		    prompt_suggestion_en = ?,
		    prompt_suggestion_ru = ?,
		    prompt_suggestion_kk = ?,
		    xp_reward = ?,
		    updated_at = NOW()
		WHERE id = ?
	`, value.LessonID, value.Position,
		value.Title.EN, value.Title.RU, value.Title.KK,
		value.Summary.EN, value.Summary.RU, value.Summary.KK,
		value.Brief.EN, value.Brief.RU, value.Brief.KK,
		value.StarterCode, successCriteriaJSON, knowledgeChecksJSON,
		value.PromptSuggestion.EN, value.PromptSuggestion.RU, value.PromptSuggestion.KK,
		value.XPReward, id,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPracticeNotFound
	}

	return nil
}

func (r *Repo) DeletePractice(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Exec(`DELETE FROM course_practices WHERE id = ?`, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPracticeNotFound
	}

	return nil
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

func (r *Repo) getPractice(ctx context.Context, condition string, args ...any) (*practicedomain.Practice, error) {
	var row practiceRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			id::text,
			lesson_id::text,
			position,
			title_en,
			title_ru,
			title_kk,
			summary_en,
			summary_ru,
			summary_kk,
			brief_en,
			brief_ru,
			brief_kk,
			starter_code,
			success_criteria,
			knowledge_checks,
			prompt_suggestion_en,
			prompt_suggestion_ru,
			prompt_suggestion_kk,
			xp_reward,
			created_at,
			updated_at
		FROM course_practices
		WHERE `+condition+`
	`, args...).Row().Scan(
		&row.ID,
		&row.LessonID,
		&row.Position,
		&row.TitleEN,
		&row.TitleRU,
		&row.TitleKK,
		&row.SummaryEN,
		&row.SummaryRU,
		&row.SummaryKK,
		&row.BriefEN,
		&row.BriefRU,
		&row.BriefKK,
		&row.StarterCode,
		&row.SuccessCriteriaJSONRaw,
		&row.KnowledgeChecksJSONRaw,
		&row.PromptSuggestionEN,
		&row.PromptSuggestionRU,
		&row.PromptSuggestionKK,
		&row.XPReward,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPracticeNotFound
		}
		return nil, err
	}

	return convertRow(row)
}

func convertRow(row practiceRow) (*practicedomain.Practice, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, err
	}
	lessonID, err := uuid.Parse(row.LessonID)
	if err != nil {
		return nil, err
	}
	successCriteria, err := unmarshalLocales(row.SuccessCriteriaJSONRaw)
	if err != nil {
		return nil, err
	}
	knowledgeChecks, err := unmarshalLocales(row.KnowledgeChecksJSONRaw)
	if err != nil {
		return nil, err
	}

	return &practicedomain.Practice{
		ID:       id,
		LessonID: lessonID,
		Position: row.Position,
		Title: practicedomain.Locale{
			EN: row.TitleEN,
			RU: row.TitleRU,
			KK: row.TitleKK,
		},
		Summary: practicedomain.Locale{
			EN: row.SummaryEN,
			RU: row.SummaryRU,
			KK: row.SummaryKK,
		},
		Brief: practicedomain.Locale{
			EN: row.BriefEN,
			RU: row.BriefRU,
			KK: row.BriefKK,
		},
		StarterCode:     row.StarterCode,
		SuccessCriteria: successCriteria,
		KnowledgeChecks: knowledgeChecks,
		PromptSuggestion: practicedomain.Locale{
			EN: row.PromptSuggestionEN,
			RU: row.PromptSuggestionRU,
			KK: row.PromptSuggestionKK,
		},
		XPReward:  row.XPReward,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func marshalLocales(values []practicedomain.Locale) (string, error) {
	if values == nil {
		values = []practicedomain.Locale{}
	}
	raw, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func unmarshalLocales(raw []byte) ([]practicedomain.Locale, error) {
	if len(raw) == 0 {
		return []practicedomain.Locale{}, nil
	}

	var values []practicedomain.Locale
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	if values == nil {
		return []practicedomain.Locale{}, nil
	}
	return values, nil
}
