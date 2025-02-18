package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"
)

// Create ...
type (
	CreateModuleInput struct {
		Module tables.Module
		_      struct{}
	}

	CreateModuleOutput struct {
		Module tables.Module
		_      struct{}
	}
)

// Get ...
type (
	GetModuleInput struct {
		ID string
		_  struct{}
	}

	GetModuleOutput struct {
		Module tables.Module
		_      struct{}
	}
)

// List ...
type (
	ListModulesWhere struct {
		ModuleCategoryID *string `where:"module-category-id"`
		State            *string `where:"state"`
		Public           *bool   `where:"public"`
		CreatorID        *string `where:"creator-id"`
		querydecoder.Query
		_ struct{}
	}

	ListModulesInput struct {
		Pagination *Pagination
		Where      ListModulesWhere
		_          struct{}
	}

	ListModulesOutput struct {
		Modules []tables.Module
		Count   int64
		_       struct{}
	}
)

// Update ...
type (
	ModuleUpdateInfo struct {
		Name        *string
		Description *string
		Location    *string
		State       *string
		Public      *bool
		_           struct{}
	}

	UpdateModuleInput struct {
		ID         string
		UpdateData *ModuleUpdateInfo
		_          struct{}
	}

	UpdateModuleOutput struct {
		Module tables.Module
		_      struct{}
	}
)

// Delete ...
type (
	DeleteModuleWhere struct {
		ID        *string
		CreatorID *string
		Public    *bool
		querydecoder.Query
		_ struct{}
	}

	DeleteModuleInput struct {
		Where DeleteModuleWhere
		_     struct{}
	}

	DeleteModuleOutput struct {
		ID []string
		_  struct{}
	}
)

// Count ...
type (
	CountModuleInput struct {
		ID string
		_  struct{}
	}

	CountModuleOutput struct {
		Count int64
		_     struct{}
	}
)
