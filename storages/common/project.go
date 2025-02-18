package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"
)

// Create ...
type (
	CreateProjectInput struct {
		Project tables.Project
		_       struct{}
	}

	CreateProjectOutput struct {
		Project tables.Project
		_       struct{}
	}
)

// Get ...
type (
	GetProjectInput struct {
		ID string
		_  struct{}
	}

	GetProjectOutput struct {
		Project tables.Project
		_       struct{}
	}
)

// List ...
type (
	ListProjectsWhere struct {
		querydecoder.Query
		_ struct{}
	}

	ListProjectsInput struct {
		Pagination *Pagination
		Where      ListProjectsWhere
		_          struct{}
	}

	ListProjectsOutput struct {
		Projects []tables.Project
		Count    int64
		_        struct{}
	}
)

// Delete ...
type (
	DeleteProjectWhere struct {
		ID *string
		querydecoder.Query
		_ struct{}
	}

	DeleteProjectInput struct {
		Where DeleteProjectWhere
		_     struct{}
	}

	DeleteProjectOutput struct {
		ID []string
		_  struct{}
	}
)
