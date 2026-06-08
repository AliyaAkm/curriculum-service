package practice

import (
	"context"
	"curriculum-service/internal/domain"
	practicedomain "curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repo) Create(ctx context.Context, value practicedomain.Task) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO practice_tasks (
			id,
			lesson_id,
			position,
			title,
			description,
			language,
			starter_code,
			expected_output,
			xp_reward,
			check_type
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, value.ID,
		value.LessonID,
		value.Position,
		value.Title,
		value.Description,
		value.Language,
		value.StarterCode,
		value.ExpectedOutput,
		value.XPReward,
		value.CheckType,
	).Error
}

func (r *Repo) Update(ctx context.Context, id uuid.UUID, value practicedomain.TaskUpdate) error {
	updates := map[string]any{
		"updated_at": gorm.Expr("NOW()"),
	}
	if value.Title != nil {
		updates["title"] = *value.Title
	}
	if value.Description != nil {
		updates["description"] = *value.Description
	}
	if value.Language != nil {
		updates["language"] = *value.Language
	}
	if value.StarterCode != nil {
		updates["starter_code"] = *value.StarterCode
	}
	if value.ExpectedOutput != nil {
		updates["expected_output"] = *value.ExpectedOutput
	}
	if value.XPReward != nil {
		updates["xp_reward"] = *value.XPReward
	}
	if value.CheckType != nil {
		updates["check_type"] = *value.CheckType
	}

	result := r.db.WithContext(ctx).
		Table("practice_tasks").
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPracticeNotFound
	}
	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Exec(`
		DELETE FROM practice_tasks
		WHERE id = ?
	`, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPracticeNotFound
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*practicedomain.Task, error) {
	var row practicedomain.Task
	if err := r.db.WithContext(ctx).Raw(baseSelect()+`
		WHERE pt.id = ?
	`, id).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.ID == uuid.Nil {
		return nil, domain.ErrPracticeNotFound
	}
	return &row, nil
}

func (r *Repo) ListByLesson(ctx context.Context, lessonID uuid.UUID) ([]practicedomain.Task, error) {
	var rows []practicedomain.Task
	err := r.db.WithContext(ctx).Raw(baseSelect()+`
		WHERE pt.lesson_id = ?
		ORDER BY pt.position ASC, pt.created_at ASC
	`, lessonID).Scan(&rows).Error
	if rows == nil {
		rows = []practicedomain.Task{}
	}
	return rows, err
}

func baseSelect() string {
	return `
		SELECT
			pt.id,
			pt.lesson_id,
			cl.module_id,
			cm.course_id,
			pt.position,
			pt.title,
			pt.description,
			pt.language,
			pt.starter_code,
			pt.expected_output,
			pt.xp_reward,
			pt.check_type,
			pt.created_at,
			pt.updated_at
		FROM practice_tasks pt
		INNER JOIN course_lessons cl ON cl.id = pt.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
	`
}
