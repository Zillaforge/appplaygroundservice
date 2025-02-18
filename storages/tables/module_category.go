package tables

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// ModuleCategory
type ModuleCategory struct {
	ID          string `gorm:"not null;primary_key;type:varchar(36)"`
	Name        string `gorm:"not null;unique;type:varchar(255)"`
	Description string `gorm:"type:text;default:NULL"`
	CreatorID   string `gorm:"not null;type:varchar(36)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	_           struct{}
}

// BeforeCreate ...
func (p *ModuleCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		for {
			p.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&ModuleCategory{}).Where("id = ?", p.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
