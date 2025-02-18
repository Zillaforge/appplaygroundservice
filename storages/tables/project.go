package tables

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// Project
type Project struct {
	ID        string `gorm:"not null;primary_key;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	_         struct{}
}

// BeforeCreate ...
func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		for {
			p.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&Project{}).Where("id = ?", p.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
