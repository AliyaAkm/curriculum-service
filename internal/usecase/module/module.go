package module

import (
	"context"
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/module"
	dtomodule "curriculum-service/internal/http/dto/module"
	"github.com/google/uuid"
)

const freePreviewModulesLimit = 2

func (u *UseCase) GetAllModules(ctx context.Context, query dtomodule.GetModuleQuery) ([]module.Module, error) {
	return u.repo.GetAllModules(ctx, query)
}

func (u *UseCase) GetAllModulesForUser(ctx context.Context, userID uuid.UUID, query dtomodule.GetModuleQuery, hasFullAccess bool) ([]module.Module, error) {
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

func (u *UseCase) GetModuleByIDForUser(ctx context.Context, userID uuid.UUID, id uuid.UUID, hasFullAccess bool) (*module.Module, error) {
	moduleEntity, err := u.repo.GetModuleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if hasFullAccess {
		return moduleEntity, nil
	}

	hasSubscription, err := u.repo.HasSubscription(ctx, userID, moduleEntity.CourseID)
	if err != nil {
		return nil, err
	}
	if hasSubscription {
		return moduleEntity, nil
	}

	allowed, err := u.repo.IsModuleInFreePreview(ctx, id, freePreviewModulesLimit)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, domain.ErrForbidden
	}

	return moduleEntity, nil
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
