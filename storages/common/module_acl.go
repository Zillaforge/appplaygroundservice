package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"
)

// Create ...
type (
	CreateModuleAclBatchInput struct {
		ModuleAcls []tables.ModuleAcl
		_          struct{}
	}

	CreateModuleAclBatchOutput struct {
		ModuleAcls []tables.ModuleAcl
		_          struct{}
	}
)

// Get ...
type (
	GetModuleAclInput struct {
		ID string
		_  struct{}
	}

	GetModuleAclOutput struct {
		ModuleAcl tables.ModuleAcl
		_         struct{}
	}
)

// List ...
type (
	ListModuleAclsWhere struct {
		ModuleID  *string `where:"module-id"`
		ProjectID *string `where:"project-id"`
		querydecoder.Query
		_ struct{}
	}

	ListModuleAclsInput struct {
		Pagination *Pagination
		Where      ListModuleAclsWhere
		_          struct{}
	}

	ListModuleAclsOutput struct {
		ModuleAcls []tables.ModuleAcl
		Count      int64
		_          struct{}
	}
)

// Delete ...
type (
	DeleteModuleAclWhere struct {
		ID        *string
		ModuleID  *string `where:"module-id"`
		ProjectID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteModuleAclInput struct {
		Where DeleteModuleAclWhere
		_     struct{}
	}

	DeleteModuleAclOutput struct {
		ID []string
		_  struct{}
	}
)
