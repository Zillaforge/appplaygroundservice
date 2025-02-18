package versions

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type id005FloatingIPAddressMigration struct {
	FloatingIPAddress string `gorm:"type:varchar(36);default:''"`
}

func getID005Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.5",
		Migrate: func(tx *gorm.DB) error {
			return id005AddFloatingIPAddressColumn(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return id005DropFloatingIPAddressColumn(tx)
		},
	}
}

// add floating_ip_address column to instance table
func id005AddFloatingIPAddressColumn(tx *gorm.DB) (err error) {
	return tx.Table("instance").Migrator().AddColumn(&id005FloatingIPAddressMigration{}, "floating_ip_address")
}

func id005DropFloatingIPAddressColumn(tx *gorm.DB) (err error) {
	return tx.Table("instance").Migrator().DropColumn(&id005FloatingIPAddressMigration{}, "floating_ip_address")
}
