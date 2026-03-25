package module

import (
	"context"
	"curriculum-service/internal/domain/module"
	dtomodule "curriculum-service/internal/http/dto/module"
	"github.com/google/uuid"
)

func (u *UseCase) GetAllModules(ctx context.Context, query dtomodule.GetModuleQuery) ([]module.Module, error) {
	return u.repo.GetAllModules(ctx, query)
}

func (u *UseCase) CreateModule(ctx context.Context, value *module.Module) (*module.Module, error) {
	value.ID = uuid.New()
	err := u.repo.CreateModule(ctx, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetModuleByID(ctx, value.ID)
}

func (u *UseCase) GetModuleByID(ctx context.Context, id uuid.UUID) (*module.Module, error) {
	return u.repo.GetModuleByID(ctx, id)
}

func (u *UseCase) UpdateModule(ctx context.Context, id uuid.UUID, value *module.Module) (*module.Module, error) {
	err := u.repo.UpdateModule(ctx, id, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetModuleByID(ctx, id)
}

func (u *UseCase) DeleteModule(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteModule(ctx, id)
}
