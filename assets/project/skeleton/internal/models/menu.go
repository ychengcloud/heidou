package models

import (
	"{{ . }}/internal/gen/models"
)

//Menu
type Menu struct {
	*models.Menu
	Children []*models.Menu `gorm:"column:children" form:"children" json:"children"`
}
