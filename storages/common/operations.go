package common

import "context"

// Operations ...
type Operations interface {
	ProjectCRUDInterface
	ModuleCategoryCRUDInterface
	ModuleCRUDInterface
	ModuleAclRUDInterface
	ApplicationCRUDInterface
	InstanceCRUDInterface
	ModuleJoinModuleAclCRUDInterface
	MeteringCRUDInterface
	AppCredentialCRUDInterface
}

type ProjectCRUDInterface interface {
	ListProjects(ctx context.Context, input *ListProjectsInput) (output *ListProjectsOutput, err error)
	CreateProject(ctx context.Context, input *CreateProjectInput) (output *CreateProjectOutput, err error)
	GetProject(ctx context.Context, input *GetProjectInput) (output *GetProjectOutput, err error)
	DeleteProject(ctx context.Context, input *DeleteProjectInput) (output *DeleteProjectOutput, err error)
}

type ModuleCategoryCRUDInterface interface {
	ListModuleCategories(ctx context.Context, input *ListModuleCategoriesInput) (output *ListModuleCategoriesOutput, err error)
	CreateModuleCategory(ctx context.Context, input *CreateModuleCategoryInput) (output *CreateModuleCategoryOutput, err error)
	GetModuleCategory(ctx context.Context, input *GetModuleCategoryInput) (output *GetModuleCategoryOutput, err error)
	UpdateModuleCategory(ctx context.Context, input *UpdateModuleCategoryInput) (output *UpdateModuleCategoryOutput, err error)
	DeleteModuleCategory(ctx context.Context, input *DeleteModuleCategoryInput) (output *DeleteModuleCategoryOutput, err error)
}

type ModuleCRUDInterface interface {
	ListModules(ctx context.Context, input *ListModulesInput) (output *ListModulesOutput, err error)
	CreateModule(ctx context.Context, input *CreateModuleInput) (output *CreateModuleOutput, err error)
	GetModule(ctx context.Context, input *GetModuleInput) (output *GetModuleOutput, err error)
	UpdateModule(ctx context.Context, input *UpdateModuleInput) (output *UpdateModuleOutput, err error)
	DeleteModule(ctx context.Context, input *DeleteModuleInput) (output *DeleteModuleOutput, err error)
	CountModule(ctx context.Context, input *CountModuleInput) (output *CountModuleOutput, err error)
}

type ModuleAclRUDInterface interface {
	ListModuleAcls(ctx context.Context, input *ListModuleAclsInput) (output *ListModuleAclsOutput, err error)
	CreateModuleAclBatch(ctx context.Context, input *CreateModuleAclBatchInput) (output *CreateModuleAclBatchOutput, err error)
	GetModuleAcl(ctx context.Context, input *GetModuleAclInput) (output *GetModuleAclOutput, err error)
	DeleteModuleAcl(ctx context.Context, input *DeleteModuleAclInput) (output *DeleteModuleAclOutput, err error)
}

type ApplicationCRUDInterface interface {
	ListApplications(ctx context.Context, input *ListApplicationsInput) (output *ListApplicationsOutput, err error)
	CreateApplication(ctx context.Context, input *CreateApplicationInput) (output *CreateApplicationOutput, err error)
	GetApplication(ctx context.Context, input *GetApplicationInput) (output *GetApplicationOutput, err error)
	UpdateApplication(ctx context.Context, input *UpdateApplicationInput) (output *UpdateApplicationOutput, err error)
	DeleteApplication(ctx context.Context, input *DeleteApplicationInput) (output *DeleteApplicationOutput, err error)
}

type InstanceCRUDInterface interface {
	ListInstances(ctx context.Context, input *ListInstancesInput) (output *ListInstancesOutput, err error)
	CreateInstance(ctx context.Context, input *CreateInstanceInput) (output *CreateInstanceOutput, err error)
	GetInstance(ctx context.Context, input *GetInstanceInput) (output *GetInstanceOutput, err error)
	UpdateInstance(ctx context.Context, input *UpdateInstanceInput) (output *UpdateInstanceOutput, err error)
	DeleteInstance(ctx context.Context, input *DeleteInstanceInput) (output *DeleteInstanceOutput, err error)
}

type ModuleJoinModuleAclCRUDInterface interface {
	ListModuleJoinModuleAcls(ctx context.Context, input *ListModuleJoinModuleAclsInput) (output *ListModuleJoinModuleAclsOutput, err error)
}

type MeteringCRUDInterface interface {
	/*
		CreateMetering 在 Metering 資料表寫入一筆 Metering Record

		errors:
		- 15000000(internal server error)
		- 15000021(metering exist)
	*/
	CreateMetering(ctx context.Context, input *CreateMeteringInput) (output *CreateMeteringOutput, err error)

	/*
		ListMeterings 回傳所有 Metering Records

		errors:
		- 15000000(internal server error)
	*/
	ListMeterings(ctx context.Context, input *ListMeteringsInput) (output *ListMeteringsOutput, err error)

	/*
		CountMetering 回傳計量資料總筆數

		errors:
		- 15000000(internal server error)
	*/
	CountMetering(ctx context.Context, input *CountMeteringInput) (output *CountMeteringOutput, err error)

	/*
		UpdateMetering 更新資料庫的 Metering Record，可更新的欄位包含 EndedAt 與 LastPublishedAt

		errors:
		- 15000000(internal server error)
		- 15000015 (application not found)
	*/
	UpdateMetering(ctx context.Context, input *UpdateMeteringInput) (output *UpdateMeteringOutput, err error)

	/*
		GetMetering 回傳指定 ID 的 Metering 資訊

		errors:
		- 15000000(internal server error)
		- 15000022(metering not found)
	*/
	GetMetering(ctx context.Context, input *GetMeteringInput) (output *GetMeteringOutput, err error)

	/*
		DeleteMetering 負責刪除指定 ID 的 Metering Record

		errors:
		- 15000000(internal server error)
	*/
	DeleteMetering(ctx context.Context, input *DeleteMeteringInput) (output *DeleteMeteringOutput, err error)
}
type AppCredentialCRUDInterface interface {
	ListAppCredentials(ctx context.Context, input *ListAppCredentialsInput) (output *ListAppCredentialsOutput, err error)
	CreateAppCredential(ctx context.Context, input *CreateAppCredentialInput) (output *CreateAppCredentialOutput, err error)
	GetAppCredential(ctx context.Context, input *GetAppCredentialInput) (output *GetAppCredentialOutput, err error)
	DeleteAppCredential(ctx context.Context, input *DeleteAppCredentialInput) (output *DeleteAppCredentialOutput, err error)
}
