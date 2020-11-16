package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	genControllers "{{ . }}/internal/gen/controllers"
	"{{ . }}/internal/gen/models"
	"{{ . }}/internal/gen/services"
)

type RoleController struct {
	*genControllers.Controller
	service services.RoleService
}

func NewRoleController(logger *zap.Logger, s services.RoleService) *RoleController {
	return &RoleController{
		Controller: &genControllers.Controller{Logger: logger},
		service:    s,
	}
}

func (roleController *RoleController) Create(c *gin.Context) {
	role := &models.Role{}
	err := c.ShouldBind(role)
	if genControllers.HandleError(c, err) {
		return
	}

	err = roleController.service.Create(c.Request.Context(), role)
	if genControllers.HandleError(c, err) {
		roleController.Logger.Error("List role error", zap.Error(err))
		return
	}

	genControllers.JsonData(c, role)
}

func (roleController *RoleController) Update(c *gin.Context) {
	role := &models.Role{}
	err := c.ShouldBind(role)
	if genControllers.HandleError(c, err) {
		return
	}

	err = roleController.service.Update(c.Request.Context(), role)
	if genControllers.HandleError(c, err) {
		roleController.Logger.Error("List role error", zap.Error(err))
		return
	}

	genControllers.JsonData(c, nil)
}

func (roleController *RoleController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if genControllers.HandleError(c, err) {
		return
	}

	err = roleController.service.Delete(c.Request.Context(), id)
	if genControllers.HandleError(c, err) {
		roleController.Logger.Error("List role error", zap.Error(err))
		return
	}

	genControllers.JsonData(c, nil)
}
