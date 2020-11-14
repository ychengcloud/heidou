package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{ . }}/gen/models"
)

type EmployeeRepository interface {
	GetByName(userName string) (employee *models.Employee, err error)
	Get(id uint32) (employee *models.Employee, err error)
}

type GormEmployeeRepository struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewGormEmployeeRepository(logger *zap.Logger, db *gorm.DB) EmployeeRepository {
	return &GormEmployeeRepository{
		logger: logger.With(zap.String("type", "EmployeeRepository")),
		db:     db,
	}
}

func (s *GormEmployeeRepository) GetByName(userName string) (employee *models.Employee, err error) {
	employee = &models.Employee{}
	err = s.db.Where("username=?", userName).First(&employee).Error
	return
}

func (s *GormEmployeeRepository) Get(id uint32) (employee *models.Employee, err error) {
	employee = &models.Employee{}

	db := s.db.Model(employee)
	db = db.Preload("Roles")
	db = db.Preload("Roles.Resources")
	db = db.Preload("Roles.Actions")
	if err = db.Where("id = ?", id).First(employee).Error; err != nil {
		return nil, errors.Wrapf(err, "get employee error[id=%d]", id)
	}
	return
}
