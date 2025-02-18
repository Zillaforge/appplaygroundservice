package versions

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type id003Metering struct {
	ApplicationID   string         `gorm:"not null;primary_key;type:varchar(36)"`
	Name            string         `gorm:"not null;type:varchar(36)"`
	ProjectID       string         `gorm:"not null;type:varchar(36)"`
	Creator         string         `gorm:"not null;type:varchar(36)"`
	Instances       datatypes.JSON `gorm:"check:json_valid(instances);default:'{}'"`
	CreatedAt       time.Time
	EndedAt         time.Time
	LastPublishedAt time.Time
	_               struct{}
}

func getID003Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.3",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable("metering") {
				return tx.Table("metering").Set("gorm:table_options", "ENGINE=InnoDB CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci").Migrator().CreateTable(&id003Metering{})
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if tx.Migrator().HasTable("metering") {
				return tx.Migrator().DropTable("metering")
			}
			return nil
		},
	}
}
