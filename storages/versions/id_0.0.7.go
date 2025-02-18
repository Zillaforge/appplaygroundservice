package versions

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type id007RollbackIndex struct {
	UserID    string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
	ProjectID string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
}

type id007MigrateIndex struct {
	UserID    string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
	ProjectID string `gorm:"not null;uniqueIndex:PU;type:varchar(36)"`
	Namespace string `gorm:"not null;uniqueIndex:PU;type:varchar(255)"`
}

func getID007Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.7",
		Migrate: func(tx *gorm.DB) error {
			return id007MigrateUniqueIndex(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return id007RollbackUniqueIndex(tx)
		},
	}
}

// add floating_ip_id column to instance table
func id007MigrateUniqueIndex(tx *gorm.DB) (err error) {
	err = tx.Table("app_credential").Migrator().DropIndex(&id007RollbackIndex{}, "PU")
	if err != nil {
		return
	}
	err = tx.Table("app_credential").Migrator().CreateIndex(&id007MigrateIndex{}, "PU")
	if err != nil {
		return
	}
	return
}

func id007RollbackUniqueIndex(tx *gorm.DB) (err error) {
	err = tx.Table("app_credential").Migrator().DropIndex(&id007MigrateIndex{}, "PU")
	if err != nil {
		return
	}
	err = tx.Table("app_credential").Migrator().CreateIndex(&id007RollbackIndex{}, "PU")
	if err != nil {
		return
	}
	return
}
