package services

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"{{ . }}/gen/models"
	"{{ . }}/gen/repositories"
	"{{ . }}/pkg/auth"
)

type RoleService interface {
	Create(c context.Context, role *models.Role) error
	Update(c context.Context, role *models.Role) error
	Delete(c context.Context, ID uint64) error
}

type DefaultRoleService struct {
	logger     *zap.Logger
	Repository repositories.RoleRepository
	auth       *auth.JWTAuth
}

func NewRoleService(logger *zap.Logger, auth *auth.JWTAuth, r repositories.RoleRepository) RoleService {
	return &DefaultRoleService{
		logger:     logger.With(zap.String("type", "DefaultRolesService")),
		Repository: r,
		auth:       auth,
	}
}

func (s *DefaultRoleService) Create(c context.Context, role *models.Role) error {

	err := s.Repository.Create(role)
	if err != nil {
		return err
	}
	err = s.loadPolicy(c, role)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultRoleService) Update(c context.Context, role *models.Role) error {
	err := s.Repository.Update(role)
	if err != nil {
		return err
	}
	err = s.loadPolicy(c, role)
	if err != nil {
		return err
	}
	return nil
}

func (s *DefaultRoleService) Delete(c context.Context, id uint64) (err error) {

	err = s.Repository.Delete(uint32(id))
	if err != nil {
		return err
	}
	roleIDStr := fmt.Sprintf("%d", id)

	ok := s.auth.Enforcer.DeletePermissionsForUser(roleIDStr)
	if !ok {
		s.logger.Info("AddPermissionForUser false")
	}
	return nil
}

func (s *DefaultRoleService) loadPolicy(c context.Context, role *models.Role) error {
	roleIDStr := fmt.Sprintf("%d", role.Id)

	s.logger.Info("DeletePermissionsForRole:",
		zap.String("roleID", roleIDStr),
	)
	ok := s.auth.Enforcer.DeletePermissionsForUser(roleIDStr)
	if !ok {
		s.logger.Info("DeletePermissionsForUser false")
	}

	for _, resource := range role.Resources {
		s.logger.Info("AddPermissionForRole:",
			zap.String("roleID", roleIDStr),
			zap.String("name", resource.Name),
			zap.String("method", resource.Method),
		)
		ok := s.auth.Enforcer.AddPermissionForUser(roleIDStr, resource.Name, resource.Method)
		if !ok {
			s.logger.Info("AddPermissionForUser false")
		}

	}

	return nil
}

