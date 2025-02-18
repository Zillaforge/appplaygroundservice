package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"

	"gorm.io/datatypes"
)

// Create ...
type (
	CreateInstanceInput struct {
		Instance tables.Instance
		_        struct{}
	}

	CreateInstanceOutput struct {
		Instance tables.Instance
		_        struct{}
	}
)

// Get ...
type (
	GetInstanceInput struct {
		ID string
		_  struct{}
	}

	GetInstanceOutput struct {
		Instance tables.Instance
		_        struct{}
	}
)

// List ...
type (
	ListInstancesWhere struct {
		ApplicationID     *string `where:"application-id"`
		ProjectID         *string `where:"project-id"`
		ReferenceID       *string `where:"reference-id"`
		FloatingIPID      *string `where:"floating-ip-id"`
		FloatingIPAddress *string `where:"floating-ip-address"`
		querydecoder.Query
		_ struct{}
	}

	ListInstancesInput struct {
		Pagination *Pagination
		Where      ListInstancesWhere
		_          struct{}
	}

	ListInstancesOutput struct {
		Instances []tables.Instance
		Count     int64
		_         struct{}
	}
)

// Update ...
type (
	InstanceUpdateInfo struct {
		Name              *string
		Extra             datatypes.JSON
		FloatingIPID      *string
		FloatingIPAddress *string
		_                 struct{}
	}

	UpdateInstanceInput struct {
		ID         string
		UpdateData *InstanceUpdateInfo
		_          struct{}
	}

	UpdateInstanceOutput struct {
		Instance tables.Instance
		_        struct{}
	}
)

// Delete ...
type (
	DeleteInstanceWhere struct {
		ID            *string
		ApplicationID *string
		ProjectID     *string
		ReferenceID   *string
		querydecoder.Query
		_ struct{}
	}

	DeleteInstanceInput struct {
		Where DeleteInstanceWhere
		_     struct{}
	}

	DeleteInstanceOutput struct {
		ID []string
		_  struct{}
	}
)
