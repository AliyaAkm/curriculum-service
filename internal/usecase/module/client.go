package module

import (
	"context"
	"curriculum-service/internal/domain/module"
	dtomodule "curriculum-service/internal/http/dto/module"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllModules(ctx context.Context, query dtomodule.GetModuleQuery) ([]module.Module, error)
	GetLimitedModulesByCourseID(ctx context.Context, courseID uuid.UUID, limit int) ([]module.Module, error)
	CreateModule(ctx context.Context, value *module.Module) error
	GetModuleByID(ctx context.Context, id uuid.UUID) (*module.Module, error)
	HasSubscription(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (bool, error)
	IsModuleInFreePreview(ctx context.Context, moduleID uuid.UUID, limit int) (bool, error)
	UpdateModule(ctx context.Context, id uuid.UUID, value *module.Module) error
	DeleteModule(ctx context.Context, id uuid.UUID) error
}
