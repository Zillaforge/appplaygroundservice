package mariadb

import (
	"AppPlaygroundService/storages/mariadb/common"
	"AppPlaygroundService/storages/mariadb/execution"
	"AppPlaygroundService/storages/mariadb/operation"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/storages/versions"
	"database/sql"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Provider ...
type Provider struct {
	conn *gorm.DB
	db   *sql.DB
	mg   *gormigrate.Gormigrate

	Exec execution.Execution
	Op   operation.Operation
}

// New create a mariadb provider and init the connection,checking the endpoint is available
func New(config *common.ConnectionConfig) (provider Provider, err error) {
	provider = Provider{}
	// setting connectionInfo to execution conf
	provider.Exec.SetConfig(config)
	// connect the database
	provider.conn, err = provider.Exec.Connect()
	if err != nil {
		return provider, err
	}
	provider.db, err = provider.conn.DB()
	if err != nil {
		return provider, err
	}
	provider.mg = gormigrate.New(provider.conn, &gormigrate.Options{TableName: tables.Migrate}, versions.Get())
	provider.Op.Set(provider.conn, provider.db)
	provider.Exec.Set(provider.mg, &provider.Op, provider.conn)

	return
}
