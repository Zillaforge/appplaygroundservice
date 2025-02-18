package tables

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// Module
type Module struct {
	ID               string `gorm:"not null;primary_key;type:varchar(36)"`
	Name             string `gorm:"not null;type:varchar(255);uniqueIndex:idx_name_module_category"`
	Description      string `gorm:"type:text;default:NULL"`
	ModuleCategoryID string `gorm:"not null;type:varchar(36);uniqueIndex:idx_name_module_category"`
	Location         string `gorm:"not null;type:text"`
	State            string `gorm:"not null;type:varchar(255)"`
	Public           bool   `gorm:"not null;default:false"`
	CreatorID        string `gorm:"not null;type:varchar(36)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ModuleCategory   ModuleCategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	_                struct{}
}

// BeforeCreate ...
func (m *Module) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		for {
			m.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&Module{}).Where("id = ?", m.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
