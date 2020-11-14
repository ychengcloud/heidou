package services

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"{{ . }}/app/models"
	"{{ . }}/app/repositories"
	gm "{{ . }}/gen/models"
	gr "{{ . }}/gen/repositories"
	"{{ . }}/pkg/auth"
)

type EmployeeService interface {
	Login(c context.Context, employee *gm.Employee) (*models.LoginTokenInfo, error)
	Create(c context.Context, employee *gm.Employee) error
	Get(c context.Context, Id uint32) (*gm.Employee, error)
	GetMenus(c context.Context, Id uint32) (menus []*models.Menu, err error)
}

type DefaultEmployeeService struct {
	logger              *zap.Logger
	employeeRepositorie repositories.EmployeeRepository
	grMenu              gr.MenuRepository
	grEmployee          gr.EmployeeRepository
	auth                *auth.JWTAuth
	rootAccount         string
}

func NewEmployeeService(v *viper.Viper, logger *zap.Logger, auth *auth.JWTAuth, employeeRepositorie repositories.EmployeeRepository, grMenu gr.MenuRepository, grEmployee gr.EmployeeRepository) EmployeeService {
	rootAccount := v.GetString("app.rootAccount")

	return &DefaultEmployeeService{
		logger:              logger.With(zap.String("type", "DefaultEmployeesService")),
		employeeRepositorie: employeeRepositorie,
		grEmployee:          grEmployee,
		grMenu:              grMenu,
		auth:                auth,
		rootAccount:         rootAccount,
	}
}

func (s *DefaultEmployeeService) Login(c context.Context, employee *gm.Employee) (*models.LoginTokenInfo, error) {
	password := employee.Password
	a, err := s.employeeRepositorie.GetByName(employee.Username)

	if err != nil {
		return nil, err
	}

	if !comparePasswords(a.Password, []byte(password)) {
		return nil, gm.ErrBadPassword
	}

	claims := make(map[string]interface{})
	claims["userId"] = a.Id
	claims["userName"] = a.Username
	claims["merchantsId"] = a.MerchantsId

	token, err := s.auth.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	tokenInfo := &models.LoginTokenInfo{
		AccessToken: token.GetAccessToken(),
		TokenType:   token.GetTokenType(),
		ExpiresAt:   token.GetExpiresAt(),
	}

	return tokenInfo, nil
}

func (s *DefaultEmployeeService) Create(c context.Context, employee *gm.Employee) error {
	if employee.Username == s.rootAccount {
		return models.ErrNameServed
	}
	s.logger.Info("create Admin", zap.String("name", employee.Username))

	err := s.grEmployee.Create(employee)
	if err != nil {
		return err
	}

	enforcer := s.auth.Enforcer
	ok, err := enforcer.DeleteRolesForUser(employee.Username)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("The user does have any roles")
	}
	for _, role := range employee.Roles {
		roleIDStr := fmt.Sprintf("%d", role.Id)
		_, err := enforcer.AddRoleForUser(employee.Username, roleIDStr)
		if err != nil {
			return err
		}
	}


	return nil
}

func (s *DefaultEmployeeService) Get(c context.Context, Id uint32) (employee *gm.Employee, err error) {

	return s.employeeRepositorie.Get(Id)
}

func toTree(menus []*gm.Menu) []*models.Menu {
	menuMap := make(map[uint32]*models.Menu)
	var list []*models.Menu
	for _, menu := range menus {
		menuMap[menu.Id] = &models.Menu{Menu: menu}
		if menu.ParentId == 0 {
			list = append(list, menuMap[menu.Id])
		}
	}

	for _, menu := range menus {
		if menu.ParentId == 0 {
			continue
		}
		if item, ok := menuMap[menu.ParentId]; ok {
			item.Children = append(item.Children, menu)
		}
	}
	fmt.Printf("%d, %#v\n", len(menus), list)
	return list
}

func (s *DefaultEmployeeService) GetMenus(c context.Context, Id uint32) ([]*models.Menu, error) {
	employee, err := s.employeeRepositorie.Get(Id)
	if err != nil {
		return nil, err
	}

	menus := make([]*gm.Menu, 0)

	// 超级管理员取全部的菜单
	if employee.Username == s.rootAccount {
		offset := 0
		limit := 20
		for {
			ms, total, err := s.grMenu.List(&gm.PaginationQuery{
				Offset: uint(offset),
				Limit:  uint(limit),
			})
			if err != nil {
				return nil, err
			}

			menus = append(menus, ms...)
			if len(ms) < limit || len(menus) == int(total) {
				break
			}
		}
	} else {
		var menuIds []uint32
		var err error
		for _, role := range employee.Roles {
			for _, resource := range role.Resources {
				menuIds = append(menuIds, resource.MenuId)
			}
		}
		menus, err = s.grMenu.BatchGet(menuIds)
		if err != nil {
			return nil, err
		}
	}

	resp := toTree(menus)
	return resp, nil
}
