package tables

import (
	"github.com/google/uuid"

	"gorm.io/gorm"
)

// ModuleAcl
type ModuleAcl struct {
	ID        string  `gorm:"not null;primary_key;type:varchar(36)"`
	ModuleID  string  `gorm:"not null;type:varchar(36);uniqueIndex:idx_mid_pid"`
	ProjectID string  `gorm:"not null;type:varchar(36);uniqueIndex:idx_mid_pid"`
	Module    Module  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Project   Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	_         struct{}
}

// BeforeCreate ...
func (m *ModuleAcl) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		for {
			m.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&ModuleAcl{}).Where("id = ?", m.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
