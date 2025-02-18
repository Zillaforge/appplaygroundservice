package tables

import (
	"time"
)

// ModuleJoinModuleAcl ...
type ModuleJoinModuleAcl struct {
	ModuleCategoryID          string
	ModuleCategoryName        string
	ModuleCategoryDescription string
	ModuleCategoryCreatorID   string
	ModuleID                  string
	ModuleName                string
	ModuleDescription         string
	Location                  string
	State                     string
	Public                    bool
	ModuleCreatorID           string
	ModuleAclID               string
	AllowProjectID            string
	ModuleCategoryCreatedAt   time.Time
	ModuleCategoryUpdatedAt   time.Time
	ModuleCreatedAt           time.Time
	ModuleUpdatedAt           time.Time
	_                         struct{}
}

// ModuleJoinModuleAclView ...
func ModuleJoinModuleAclView() string {
	return `CREATE VIEW module_join_module_acl AS 
	SELECT
		module_category.id AS module_category_id,
		module_category.name AS module_category_name,
		module_category.description AS module_category_description,
		module_category.creator_id AS module_category_creator_id,

		module.id AS module_id,
		module.name AS module_name,
		module.description AS module_description,
		module.location AS location,
		module.state AS state,
		module.public AS public,
		module.creator_id AS module_creator_id,

		module_acl.id AS module_acl_id,
		module_acl.project_id AS allow_project_id,

		module_category.created_at AS module_category_created_at,
		module_category.updated_at AS module_category_updated_at,
		module.created_at AS module_created_at,
		module.updated_at AS module_updated_at
	FROM
		aps.module_category
	LEFT JOIN aps.module ON
		aps.module_category.id = aps.module.module_category_id
	LEFT JOIN aps.module_acl ON
		aps.module.id = aps.module_acl.module_id`
}
