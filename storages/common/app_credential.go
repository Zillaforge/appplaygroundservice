package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"
)

// Create ...
type (
	CreateAppCredentialInput struct {
		AppCredential tables.AppCredential
		_             struct{}
	}
	CreateAppCredentialOutput struct {
		AppCredential tables.AppCredential
		_             struct{}
	}
)

// Get ...
type (
	GetAppCredentialInput struct {
		ID        *string
		UserID    *string
		ProjectID *string
		Namespace *string
		_         struct{}
	}
	GetAppCredentialOutput struct {
		AppCredential tables.AppCredential
		_             struct{}
	}
)

// List ...
type (
	ListAppCredentialsInput struct {
		Pagination *Pagination
		Where      AppCredentialWhere
		_          struct{}
	}
	ListAppCredentialsOutput struct {
		AppCredentials []tables.AppCredential
		Count          int64
		_              struct{}
	}
)

// Delete ...
type (
	DeleteAppCredentialInput struct {
		Where AppCredentialWhere
		_     struct{}
	}
	DeleteAppCredentialOutput struct {
		ID []string
		_  struct{}
	}
)

// Common ...
type (
	AppCredentialWhere struct {
		ID        *string `where:"id"`
		UserID    *string `where:"user-id"`
		ProjectID *string `where:"project-id"`
		Namespace *string `where:"namespace"`
		querydecoder.Query
		_ struct{}
	}
)
