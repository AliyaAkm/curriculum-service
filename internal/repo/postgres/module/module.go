package module

import (
	"context"
	"curriculum-service/internal/domain/module"
	dtomodule "curriculum-service/internal/http/dto/module"
	"github.com/google/uuid"
	"strings"
)

func (r *Repo) GetAllModules(ctx context.Context, query dtomodule.GetModuleQuery) ([]module.Module, error) {
	var modules []module.Module
	db := r.db.WithContext(ctx)

	locale := strings.TrimSpace(strings.ToLower(query.Locale))
	if locale != "" {
		db = db.Where("locale = ?", locale)
	}
	if query.CourseID != uuid.Nil {
		db = db.Where("course_id = ?", query.CourseID)
	}

	db = db.Order("position ASC").Order("locale ASC")
	if query.CourseID == uuid.Nil || query.Limit > 0 || query.Page > 0 {
		limit := 10
		if query.Limit > 0 {
			limit = query.Limit
		}
		page := 1
		if query.Page > 0 {
			page = query.Page
		}
		offset := (page - 1) * limit
		db = db.Limit(limit).Offset(offset)
	}

	err := db.Find(&modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *Repo) GetModuleByCourseID(ctx context.Context, courseID uuid.UUID) ([]module.Module, error) {
	var modules []module.Module
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Order("position ASC").Find(&modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *Repo) GetLimitedModulesByCourseID(ctx context.Context, courseID uuid.UUID, limit int) ([]module.Module, error) {
	var modules []module.Module
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Order("position ASC").Limit(limit).Find(&modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *Repo) CreateModule(ctx context.Context, value *module.Module) error {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) UpdateModule(ctx context.Context, id uuid.UUID, value *module.Module) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(value).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetModuleByID(ctx context.Context, id uuid.UUID) (*module.Module, error) {
	var entity module.Module
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &entity, nil
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

func (r *Repo) IsModuleInFreePreview(ctx context.Context, moduleID uuid.UUID, limit int) (bool, error) {
	var allowed bool
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT EXISTS (
				SELECT 1
				FROM course_modules
				WHERE id = ?
				  AND position <= ?
			)
		`, moduleID, limit).
		Scan(&allowed).Error
	if err != nil {
		return false, err
	}
	return allowed, nil
}

func (r *Repo) DeleteModule(ctx context.Context, id uuid.UUID) error {
	var entity module.Module
	err := r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
