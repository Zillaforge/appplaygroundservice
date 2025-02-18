package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"

	"gorm.io/datatypes"
)

// Create ...
type (
	CreateApplicationInput struct {
		Application tables.Application
		_           struct{}
	}

	CreateApplicationOutput struct {
		Application tables.Application
		_           struct{}
	}
)

// Get ...
type (
	GetApplicationInput struct {
		ID string
		_  struct{}
	}

	GetApplicationOutput struct {
		Application tables.Application
		_           struct{}
	}
)

// List ...
type (
	ListApplicationsWhere struct {
		ModuleID  *string `where:"module-id"`
		State     *string `where:"state"`
		Namespace *string `where:"namespace"`
		Shiftable *bool   `where:"shiftable"`
		ProjectID *string `where:"project-id"`
		CreatorID *string `where:"creator-id"`
		UpdaterID *string `where:"updater-id"`
		querydecoder.Query
		_ struct{}
	}

	ListApplicationsInput struct {
		Pagination *Pagination
		Where      ListApplicationsWhere
		_          struct{}
	}

	ListApplicationsOutput struct {
		Applications []tables.Application
		Count        int64
		_            struct{}
	}
)

// Update ...
type (
	ApplicationUpdateInfo struct {
		Name        *string
		Description *string
		State       *string
		Answers     datatypes.JSON
		Namespace   *string
		Shiftable   *bool
		Extra       datatypes.JSON
		UpdaterID   string
		_           struct{}
	}

	UpdateApplicationInput struct {
		ID         string
		UpdateData *ApplicationUpdateInfo
		_          struct{}
	}

	UpdateApplicationOutput struct {
		Application tables.Application
		_           struct{}
	}
)

// Delete ...
type (
	DeleteApplicationWhere struct {
		ID        *string
		ModuleID  *string
		State     *string
		ProjectID *string
		CreatorID *string
		Shiftable *bool
		querydecoder.Query
		_ struct{}
	}

	DeleteApplicationInput struct {
		Where DeleteApplicationWhere
		_     struct{}
	}

	DeleteApplicationOutput struct {
		ID []string
		_  struct{}
	}
)
