package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"
)

// Create ...
type (
	CreateModuleCategoryInput struct {
		ModuleCategory tables.ModuleCategory
		_              struct{}
	}

	CreateModuleCategoryOutput struct {
		ModuleCategory tables.ModuleCategory
		_              struct{}
	}
)

// Get ...
type (
	GetModuleCategoryInput struct {
		ID string
		_  struct{}
	}

	GetModuleCategoryOutput struct {
		ModuleCategory tables.ModuleCategory
		_              struct{}
	}
)

// List ...
type (
	ListModuleCategoriesWhere struct {
		CreatorID *string `where:"creator-id"`
		querydecoder.Query
		_ struct{}
	}

	ListModuleCategoriesInput struct {
		Pagination     *Pagination
		Where          ListModuleCategoriesWhere
		AllowProjectID *string
		_              struct{}
	}

	ListModuleCategoriesOutput struct {
		ModuleCategories []tables.ModuleCategory
		Count            int64
		_                struct{}
	}
)

// Update ...
type (
	ModuleCategoryUpdateInfo struct {
		Description *string
		_           struct{}
	}

	UpdateModuleCategoryInput struct {
		ID         string
		UpdateData *ModuleCategoryUpdateInfo
		_          struct{}
	}

	UpdateModuleCategoryOutput struct {
		ModuleCategory tables.ModuleCategory
		_              struct{}
	}
)

// Delete ...
type (
	DeleteModuleCategoryWhere struct {
		ID        *string
		CreatorID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteModuleCategoryInput struct {
		Where DeleteModuleCategoryWhere
		_     struct{}
	}

	DeleteModuleCategoryOutput struct {
		ID []string
		_  struct{}
	}
)
