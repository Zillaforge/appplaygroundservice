package common

import (
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility/querydecoder"
)

// List ...
type (
	ListModuleJoinModuleAclsWhere struct {
		ModuleCategoryID        *string `where:"module-category-id"`
		ModuleCategoryCreatorID *string `where:"module-category-creator-id"`
		ModuleID                *string `where:"module-id"`
		State                   *string `where:"state"`
		Public                  *bool   `where:"public"`
		ModuleCreatorID         *string `where:"module-creator-id"`
		AllowProjectID          *string `where:"allow-project-id"`
		querydecoder.Query
		_ struct{}
	}

	ListModuleJoinModuleAclsInput struct {
		Pagination *Pagination
		Where      ListModuleJoinModuleAclsWhere
		ProjectID  *string
		_          struct{}
	}

	ListModuleJoinModuleAclsOutput struct {
		ModuleJoinModuleAcls []tables.ModuleJoinModuleAcl
		Count                int64
		_                    struct{}
	}
)
