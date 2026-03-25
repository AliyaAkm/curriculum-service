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
	if locale == "" {
		locale = "en"
	}
	db = db.Where("locale = ?", locale)
	if query.CourseID != uuid.Nil {
		db = db.Where("course_id = ?", query.CourseID)
	}
	limit := 10
	if query.Limit > 0 {
		limit = query.Limit
	}
	page := 1
	if query.Page > 0 {
		page = query.Page
	}
	offset := (page - 1) * limit
	err := db.Limit(limit).Offset(offset).Find(&modules).Error
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

func (r *Repo) DeleteModule(ctx context.Context, id uuid.UUID) error {
	var entity module.Module
	err := r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
