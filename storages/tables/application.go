package tables

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Application
type Application struct {
	ID          string         `gorm:"not null;primary_key;type:varchar(36)"`
	Name        string         `gorm:"not null;type:varchar(255)"`
	Description string         `gorm:"type:text;default:NULL"`
	ModuleID    string         `gorm:"not null;type:varchar(36)"`
	State       string         `gorm:"not null;type:varchar(255)"`
	Answers     datatypes.JSON `gorm:"check:json_valid(answers);default:'{}'"`
	Namespace   string         `gorm:"not null;type:varchar(255)"`
	Shiftable   bool           `gorm:"not null;default:false"`
	ProjectID   string         `gorm:"not null;type:varchar(36)"`
	CreatorID   string         `gorm:"not null;type:varchar(36)"`
	UpdaterID   string         `gorm:"nullable;type:varchar(36)"`
	Extra       datatypes.JSON `gorm:"check:json_valid(extra);default:'{}'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Module      Module  `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Project     Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	_           struct{}
}

// BeforeCreate ...
func (a *Application) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		for {
			a.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&Application{}).Where("id = ?", a.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
