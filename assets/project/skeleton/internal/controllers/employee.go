package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	gc "{{ . }}/internal/gen/controllers"
	"{{ . }}/internal/gen/models"
	"{{ . }}/internal/services"
)

//LoginRequest ...
type LoginRequest struct {
	UserName    string `form:"username" json:"username" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	CaptchaID   string `json:"captchaId" binding:"" swaggo:"true,验证码ID"`
	CaptchaCode string `json:"captchaCode" binding:"" swaggo:"true,验证码"`
}

type EmployeeController struct {
	*gc.Controller
	service services.EmployeeService
}

func NewEmployeeController(controller *gc.Controller, s services.EmployeeService) *EmployeeController {
	return &EmployeeController{
		Controller: controller,
		service:    s,
	}
}

func (ctl *EmployeeController) Signin(c *gin.Context) {

	request := &LoginRequest{}
	err := c.ShouldBind(request)
	if gc.HandleError(c, err) {
		return
	}

	employee := &models.Employee{
		Username: request.UserName,
		Password: request.Password,
	}

	token, err := ctl.service.Login(c.Request.Context(), employee)
	if gc.HandleError(c, err) {
		ctl.Logger.Error("Employee login error", zap.Error(err))
		return
	}

	gc.JsonData(c, token)
}

// Logout 登出
func (ctl *EmployeeController) Signout(c *gin.Context) {
	gc.JsonData(c, nil)
}

//GetCurrentEmployee 获取当前用户信息
func (ctl *EmployeeController) GetCurrentEmployee(c *gin.Context) {
	id, _, _, err := parseClaims(c, ctl.ClaimsKey)
	if gc.HandleError(c, err) {
		ctl.Logger.Error("parseClaims", zap.Error(err))
		return
	}

	employee, err := ctl.service.Get(c.Request.Context(), id)
	if gc.HandleError(c, err) {
		ctl.Logger.Error("Get employee error", zap.Error(err))
		return
	}
	gc.JsonData(c, employee)
}

//GetCurrentMenus 获取当前用户信息
func (ctl *EmployeeController) GetCurrentMenus(c *gin.Context) {
	id, _, _, err := parseClaims(c, ctl.ClaimsKey)
	if gc.HandleError(c, err) {
		ctl.Logger.Error("parseClaims", zap.Error(err))
		return
	}

	employee, err := ctl.service.GetMenus(c.Request.Context(), id)
	if gc.HandleError(c, err) {
		ctl.Logger.Error("Get employee error", zap.Error(err))
		return
	}
	gc.JsonData(c, employee)
}

func (ctl *EmployeeController) Create(c *gin.Context) {
	employee := &models.Employee{}
	err := c.ShouldBind(employee)
	if gc.HandleError(c, err) {
		return
	}

	err = ctl.service.Create(c.Request.Context(), employee)
	if gc.HandleError(c, err) {
		ctl.Logger.Error("List employee error", zap.Error(err))
		return
	}

	gc.JsonData(c, employee)
}