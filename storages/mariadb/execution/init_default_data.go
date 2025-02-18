package execution

import (
	"AppPlaygroundService/storages/mariadb/operation"
)

// InitDefaultData ...
func (e *Execution) InitDefaultData() (err error) {
	// set connection to operation
	e.op = &operation.Operation{}
	e.op.Set(e.conn, e.db)
	return
}
