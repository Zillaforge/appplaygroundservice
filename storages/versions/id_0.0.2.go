package versions

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type id002FloatingIPIDMigration struct {
	FloatingIPID string `gorm:"type:varchar(36);default:''"`
}

func getID002Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.2",
		Migrate: func(tx *gorm.DB) error {
			return id002AddFloatingIPIDColumn(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return id002DropFloatingIPIDColumn(tx)
		},
	}
}

// add floating_ip_id column to instance table
func id002AddFloatingIPIDColumn(tx *gorm.DB) (err error) {
	return tx.Table("instance").Migrator().AddColumn(&id002FloatingIPIDMigration{}, "floating_ip_id")
}

func id002DropFloatingIPIDColumn(tx *gorm.DB) (err error) {
	return tx.Table("instance").Migrator().DropColumn(&id002FloatingIPIDMigration{}, "floating_ip_id")
}
