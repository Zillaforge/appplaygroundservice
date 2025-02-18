package tables

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Instance
type Instance struct {
	ID                string         `gorm:"not null;primary_key;type:varchar(36)"`
	Name              string         `gorm:"not null;type:varchar(255)"`
	ApplicationID     string         `gorm:"not null;type:varchar(36)"`
	ProjectID         string         `gorm:"not null;type:varchar(36)"`
	ReferenceID       string         `gorm:"not null;type:varchar(255)"`
	Extra             datatypes.JSON `gorm:"check:json_valid(extra);default:'{}'"`
	FloatingIPID      string         `gorm:"type:varchar(36);default:''"`
	FloatingIPAddress string         `gorm:"type:varchar(36);default:''"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Application       Application `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Project           Project     `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	_                 struct{}
}

// BeforeCreate ...
func (i *Instance) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		for {
			i.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&Instance{}).Where("id = ?", i.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
