package versions

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type id006Migration struct {
	Namespace string `gorm:"not null;type:varchar(255)"`
}

func getID006Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.6",
		Migrate: func(tx *gorm.DB) error {
			return id006AddNamespaceColumn(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return id006DropNamespaceColumn(tx)
		},
	}
}

// add floating_ip_id column to instance table
func id006AddNamespaceColumn(tx *gorm.DB) (err error) {
	return tx.Table("app_credential").Migrator().AddColumn(&id006Migration{}, "namespace")
}

func id006DropNamespaceColumn(tx *gorm.DB) (err error) {
	return tx.Table("app_credential").Migrator().DropColumn(&id006Migration{}, "namespace")
}
