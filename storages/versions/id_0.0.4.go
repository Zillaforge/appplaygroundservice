package versions

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type id004AppCredential struct {
	UserID    string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
	ProjectID string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
	ID        string `gorm:"not null;type:varchar(36)"`
	Name      string `gorm:"not null;type:varchar(255)"`
	Secret    string `gorm:"not null;type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	_         struct{}
}

func getID004Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.4",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable("app_credential") {
				return tx.Table("app_credential").Set("gorm:table_options", "ENGINE=InnoDB CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci").Migrator().CreateTable(&id004AppCredential{})
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if tx.Migrator().HasTable("app_credential") {
				return tx.Migrator().DropTable("app_credential")
			}
			return nil
		},
	}
}
