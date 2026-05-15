package lesson

import (
	"context"
	"curriculum-service/internal/domain/keypoint"
	lessondomain "curriculum-service/internal/domain/lesson"
	"curriculum-service/internal/domain/outcome"
	"curriculum-service/internal/domain/summary"
	"curriculum-service/internal/domain/theorycontent"
	"curriculum-service/internal/domain/title"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repo) GetAllLessons(ctx context.Context, moduleID uuid.UUID) ([]lessondomain.LessonModel, error) {
	var rows []lessondomain.LessonModel

	err := r.db.WithContext(ctx).
		Where("module_id = ?", moduleID).
		Preload("Titles.Locale").
		Preload("Summaries.Locale").
		Preload("Outcomes.Locale").
		Preload("TheoryContents.Locale").
		Preload("KeyPoints.Locale").
		Order("position ASC").
		Order("created_at ASC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("get all lessons: %w", err)
	}

	return rows, nil
}

func (r *Repo) GetLessonByID(ctx context.Context, id uuid.UUID) (*lessondomain.LessonModel, error) {
	var lessonEntity lessondomain.LessonModel

	err := r.db.WithContext(ctx).
		Preload("Titles.Locale").
		Preload("Summaries.Locale").
		Preload("Outcomes.Locale").
		Preload("TheoryContents.Locale").
		Preload("KeyPoints.Locale").
		Where("id = ?", id).
		First(&lessonEntity).Error
	if err != nil {
		return nil, err
	}

	return &lessonEntity, nil
}

func (r *Repo) GetCourseIDByModuleID(ctx context.Context, moduleID uuid.UUID) (uuid.UUID, error) {
	var courseID string
	err := r.db.WithContext(ctx).
		Table("course_modules").
		Select("course_id").
		Where("id = ?", moduleID).
		Scan(&courseID).Error
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(courseID)
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

func (r *Repo) CreateLesson(ctx context.Context, value *lessondomain.LessonModel) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		lesson := lessondomain.LessonModel{
			ID:              value.ID,
			ModuleID:        value.ModuleID,
			Position:        value.Position,
			DurationMinutes: value.DurationMinutes,
			XPReward:        value.XPReward,
			CodeSnippet:     value.CodeSnippet,
			ExampleOutput:   value.ExampleOutput,
			VideoObjectKey:  value.VideoObjectKey,
		}

		if err := tx.Create(&lesson).Error; err != nil {
			return err
		}

		for i := range value.Titles {
			value.Titles[i].LessonID = value.ID
		}
		for i := range value.Summaries {
			value.Summaries[i].LessonID = value.ID
		}
		for i := range value.Outcomes {
			value.Outcomes[i].LessonID = value.ID
		}
		for i := range value.TheoryContents {
			value.TheoryContents[i].LessonID = value.ID
		}
		for i := range value.KeyPoints {
			value.KeyPoints[i].LessonID = value.ID
		}

		if len(value.Titles) > 0 {
			if err := tx.Create(&value.Titles).Error; err != nil {
				return err
			}
		}
		if len(value.Summaries) > 0 {
			if err := tx.Create(&value.Summaries).Error; err != nil {
				return err
			}
		}
		if len(value.Outcomes) > 0 {
			if err := tx.Create(&value.Outcomes).Error; err != nil {
				return err
			}
		}
		if len(value.TheoryContents) > 0 {
			if err := tx.Create(&value.TheoryContents).Error; err != nil {
				return err
			}
		}
		if len(value.KeyPoints) > 0 {
			if err := tx.Create(&value.KeyPoints).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
func (r *Repo) UpdateLesson(ctx context.Context, id uuid.UUID, value *lessondomain.LessonModel) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Model(&lessondomain.LessonModel{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"module_id":        value.ModuleID,
				"position":         value.Position,
				"duration_minutes": value.DurationMinutes,
				"xp_reward":        value.XPReward,
				"code_snippet":     value.CodeSnippet,
				"example_output":   value.ExampleOutput,
				"video_object_key": value.VideoObjectKey,
			}).Error
		if err != nil {
			return err
		}
		if err = r.deleteLessonsAttributeTx(ctx, tx, id); err != nil {
			return err
		}
		for i := range value.Titles {
			value.Titles[i].LessonID = id
		}
		for i := range value.Summaries {
			value.Summaries[i].LessonID = id
		}
		for i := range value.Outcomes {
			value.Outcomes[i].LessonID = id
		}
		for i := range value.TheoryContents {
			value.TheoryContents[i].LessonID = id
		}
		for i := range value.KeyPoints {
			value.KeyPoints[i].LessonID = id
		}

		if err = r.createLessonsAttributeTx(ctx, tx, value); err != nil {
			return err
		}
		return nil
	})
}

func (r *Repo) UpdateLessonVideoObjectKey(ctx context.Context, id uuid.UUID, videoObjectKey *string) error {
	return r.db.WithContext(ctx).
		Model(&lessondomain.LessonModel{}).
		Where("id = ?", id).
		Update("video_object_key", videoObjectKey).
		Error
}

func (r *Repo) createLessonsAttributeTx(ctx context.Context, tx *gorm.DB, value *lessondomain.LessonModel) error {
	if len(value.Titles) > 0 {
		if err := tx.WithContext(ctx).Create(&value.Titles).Error; err != nil {
			return err
		}
	}
	if len(value.Summaries) > 0 {
		if err := tx.WithContext(ctx).Create(&value.Summaries).Error; err != nil {
			return err
		}
	}
	if len(value.Outcomes) > 0 {
		if err := tx.WithContext(ctx).Create(&value.Outcomes).Error; err != nil {
			return err
		}
	}
	if len(value.TheoryContents) > 0 {
		if err := tx.WithContext(ctx).Create(&value.TheoryContents).Error; err != nil {
			return err
		}
	}
	if len(value.KeyPoints) > 0 {
		if err := tx.WithContext(ctx).Create(&value.KeyPoints).Error; err != nil {
			return err
		}
	}
	return nil
}
func (r *Repo) deleteLessonsAttributeTx(ctx context.Context, tx *gorm.DB, lessonID uuid.UUID) error {
	var err error

	if err = tx.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		Delete(&title.LessonTitleModel{}).Error; err != nil {
		return err
	}

	if err = tx.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		Delete(&summary.LessonSummaryModel{}).Error; err != nil {
		return err
	}

	if err = tx.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		Delete(&outcome.LessonOutcomeModel{}).Error; err != nil {
		return err
	}

	if err = tx.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		Delete(&theorycontent.LessonTheoryContentModel{}).Error; err != nil {
		return err
	}

	if err = tx.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		Delete(&keypoint.LessonKeyPointModel{}).Error; err != nil {
		return err
	}

	return nil
}
