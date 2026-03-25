package module

import (
	"context"
	"curriculum-service/internal/domain/module"
	dtomodule "curriculum-service/internal/http/dto/module"
	"github.com/google/uuid"
)

type client interface {
	GetAllModules(ctx context.Context, query dtomodule.GetModuleQuery) ([]module.Module, error)
	CreateModule(ctx context.Context, value *module.Module) (*module.Module, error)
	GetModuleByID(ctx context.Context, id uuid.UUID) (*module.Module, error)
	UpdateModule(ctx context.Context, id uuid.UUID, value *module.Module) (*module.Module, error)
	DeleteModule(ctx context.Context, id uuid.UUID) error
}
