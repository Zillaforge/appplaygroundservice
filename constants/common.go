package constants

const (
	Name           = "AppPlaygroundService"
	PascalCaseName = "AppPlaygroundService"
	SnakeCaseName  = "app_playground_service"
	KebabCaseName  = "app-playground-service"
	UpperAbbrName  = "APS"
	LowerAbbrName  = "aps"

	Kind                 = PascalCaseName
	Version              = "0.0.5"
	APIPrefix            = "/" + LowerAbbrName + "/api/"
	APIVersion           = "v1"
	GlobalConfigPath     = "etc/ASUS"
	GlobalConfigFilename = KebabCaseName + ".yaml"
	ProductUUIDFilePath  = "/sys/class/dmi/id/product_uuid"

	//Workflow log
	RequestID      = "Request-ID"
	Middleware     = "Middleware"
	Authentication = "Authentication"
	Storage        = "Storage"
	GRPC           = "GRPC"
	Controller     = "Controller"
	Server         = "Server"
	Module         = "Module"
	Cmd            = "Cmd"
	FSM            = "FSM"
	Plugin         = "Plugin"
	EventConsume   = "EventConsume"
	Task           = "Task"
	OpskResource   = "OpskResource"

	// Query Keys ...
	QueryToken = "token"

	// Header Keys ...
	HdrHostID                        = "Host-ID"
	HdrLocationID                    = "Location-ID"
	HdrVersionID                     = "Version-ID"
	HdrAuthorization                 = "Authorization"
	HdrProjectIDFromKong             = "Project-ID"
	HdrUserIDFromKong                = "User-ID"
	HdrUserRoleFromKong              = "User-Role"
	HdrProjectActiveFromKong         = "Project-Active"
	HdrSystemAdminFromKong           = "System-Admin"
	HdrUserAccountFromKong           = "User-Account"
	HdrSAATUserIDFromKong            = "SAAT-User-ID"
	HdrNamespace                     = "X-Namespace"
	HdrContentLanguage               = "Content-Language"
	HdrLanguage                      = "X-Language"
	HdrProjectExtraResourceReviewAPS = "Project-Extra-resourceReview-aps"

	// Context
	CtxLocationID       = HdrLocationID
	CtxHostID           = HdrHostID
	CtxUserID           = "ctxUserID"
	CtxUserAccount      = "ctxUserAccount"
	CtxProjectID        = "ctxProjectID"
	CtxTenantRole       = "ctxTenantRole"
	CtxOperationName    = "ctxOperationName"
	CtxCreator          = "ctxCreator"
	CtxSAATUserID       = "ctxSAATUserID"
	CtxModuleCategoryID = "ctxModuleCategoryID"
	CtxModuleID         = "ctxModuleID"
	CtxApplicationID    = "ctxApplicationID"
	CtxNamespace        = "ctxNamespace"
	CtxInstanceID       = "ctxInstanceID"
	CtxLanguage         = "ctxLanguage"
	CtxResourceReview   = "ctxResourceReview"

	// Params
	ParamProjectID        = "project-id"
	ParamModuleCategoryID = "module-category-id"
	ParamModuleID         = "module-id"
	ParamApplicationID    = "application-id"
	ParamInstanceID       = "instance-id"

	//ReconcileKey ...
	ReconcileKey = "eventpublish"
	SyncKey      = "syncevent"
	AsyncKey     = "asyncevent"

	// FSM Task
	Application = "application"
)

type CtxWithExtraVal struct{}
