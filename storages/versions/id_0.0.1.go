package versions

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type id001ShiftableMigration struct {
	Shiftable bool `gorm:"not null;default:false"`
}

func getID001Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.1",
		Migrate: func(tx *gorm.DB) error {
			return id001AddShiftableColumn(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return id001DropShiftableColumn(tx)
		},
	}
}

// add shiftable column to application table
func id001AddShiftableColumn(tx *gorm.DB) (err error) {
	return tx.Table("application").Migrator().AddColumn(&id001ShiftableMigration{}, "shiftable")
}

func id001DropShiftableColumn(tx *gorm.DB) (err error) {
	return tx.Table("application").Migrator().DropColumn(&id001ShiftableMigration{}, "shiftable")
}
