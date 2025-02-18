package versions

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type id008Migration struct {
	Extra datatypes.JSON `gorm:"check:json_valid(extra);default:'{}'"`
}

func getID008Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.8",
		Migrate: func(tx *gorm.DB) error {
			return id008AddExtraColumn(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return id008DropExtraColumn(tx)
		},
	}
}

// add extra column to application table
func id008AddExtraColumn(tx *gorm.DB) (err error) {
	return tx.Table("application").Migrator().AddColumn(&id008Migration{}, "extra")
}

func id008DropExtraColumn(tx *gorm.DB) (err error) {
	return tx.Table("application").Migrator().DropColumn(&id008Migration{}, "extra")
}
