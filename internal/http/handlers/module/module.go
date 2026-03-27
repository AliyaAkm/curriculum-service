package module

import (
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/module"
	dtomodule "curriculum-service/internal/http/dto/module"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) GetAllModules(c *gin.Context) {
	var query dtomodule.GetModuleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid query params")
		return
	}
	result, err := h.client.GetAllModules(c.Request.Context(), query)
	if err != nil {
		writeModuleError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertModules(result))
}

func (h *Handler) CreateModule(c *gin.Context) {
	request := dtomodule.ModuleRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.CreateModule(c.Request.Context(), convertModuleRequest(request))
	if err != nil {
		writeModuleError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertModule(result))
}

func (h *Handler) GetModuleByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid module id")
		return
	}

	result, err := h.client.GetModuleByID(c.Request.Context(), id)
	if err != nil {
		writeModuleError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertModule(result))
}

func (h *Handler) DeleteModule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid module id")
		return
	}

	err = h.client.DeleteModule(c.Request.Context(), id)
	if err != nil {
		writeModuleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) UpdateModule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid module id")
		return
	}

	request := dtomodule.ModuleRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.UpdateModule(c.Request.Context(), id, convertModuleRequest(request))
	if err != nil {
		writeModuleError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertModule(result))
}

func writeModuleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request query")
	default:
		c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertModules(resp []module.Module) []dtomodule.Modules {
	modules := make([]dtomodule.Modules, len(resp))
	for i := range resp {
		modules[i] = dtomodule.Modules{
			ID:        resp[i].ID,
			CourseID:  resp[i].CourseID,
			Title:     resp[i].Title,
			Summary:   resp[i].Description,
			Locale:    resp[i].Locale,
			CreatedAt: resp[i].CreatedAt,
			UpdatedAt: resp[i].UpdatedAt,
		}
	}
	return modules
}
func convertModule(resp *module.Module) dtomodule.Modules {
	return dtomodule.Modules{
		ID:        resp.ID,
		CourseID:  resp.CourseID,
		Title:     resp.Title,
		Summary:   resp.Description,
		Locale:    resp.Locale,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}

func convertModuleRequest(resp dtomodule.ModuleRequest) *module.Module {
	locale := resp.Locale
	if locale == "" {
		locale = "en"
	}

	return &module.Module{
		CourseID:    resp.CourseID,
		Title:       resp.Title,
		Description: resp.Summary,
		Locale:      locale,
	}
}
