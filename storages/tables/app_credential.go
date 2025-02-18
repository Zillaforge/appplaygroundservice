package tables

import "time"

type AppCredential struct {
	UserID    string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
	ProjectID string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
	ID        string `gorm:"not null;type:varchar(36)"`
	Name      string `gorm:"not null;type:varchar(255)"`
	Secret    string `gorm:"not null;type:text"`
	Namespace string `gorm:"not null;uniqueIndex:PU;type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	_         struct{}
}
