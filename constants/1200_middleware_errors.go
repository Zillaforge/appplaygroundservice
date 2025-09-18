package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1200xxxx: Middleware

	MidInternalServerErrorErrCode       = 12000000
	MidInternalServerErrorErrMsg        = "internal server error"
	MidPermissionDeniedErrCode          = 12000001
	MidPermissionDeniedErrMsg           = "permission denied"
	MidUserHasBeenFrozenErrCode         = 12000002
	MidUserHasBeenFrozenErrMsg          = "the user has been frozen, please contact administrator"
	MidMembershipHasBeenFrozenErrCode   = 12000003
	MidMembershipHasBeenFrozenErrMsg    = "the membership has been frozen, please contact tenant of admin"
	MidProjectNotFoundErrCode           = 12000004
	MidProjectNotFoundErrMsg            = "project (%s) not found"
	MidIncorrectFormatErrCode           = 12000005
	MidIncorrectFormatErrMsg            = "incorrect format of authentication"
	MidModuleCategoryNotFoundErrCode    = 12000006
	MidModuleCategoryNotFoundErrCodeMsg = "module category (%s) not found"
	MidModuleNotFoundErrCode            = 12000007
	MidModuleNotFoundErrCodeMsg         = "module (%s) not found"
	MidModuleIsReadOnlyErrCode          = 12000008
	MidModuleIsReadOnlyErrMsg           = "module (%s) is read only for users of other projects"
	MidApplicationNotFoundErrCode       = 12000009
	MidApplicationNotFoundErrMsg        = "application (%s) not found"
	MidNamespaceNotAllowErrCode         = 12000010
	MidNamespaceNotAllowErrMsg          = "namespace not allow"
	MidApplicationIsReadOnlyErrCode     = 12000011
	MidApplicationIsReadOnlyErrMsg      = "application (%s) is read only for users of other projects"
	MidInstanceNotFoundErrCode          = 12000012
	MidInstanceNotFoundErrMsg           = "instance (%s) not found"
	MidInstanceIsReadOnlyErrCode        = 12000013
	MidInstanceIsReadOnlyErrMsg         = "instance (%s) is read only for users of other projects"
)

var (
	// 1200xxxx: Middleware

	// 12000000(internal server error)
	MidInternalServerErrorErr = tkErr.Error(MidInternalServerErrorErrCode, MidInternalServerErrorErrMsg)
	// 12000001(permission denied)
	MidPermissionDeniedErr = tkErr.Error(MidPermissionDeniedErrCode, MidPermissionDeniedErrMsg)
	// 12000002(the user has been frozen, please contact administrator)
	MidUserHasBeenFrozenErr = tkErr.Error(MidUserHasBeenFrozenErrCode, MidUserHasBeenFrozenErrMsg)
	// 12000003(the membership has been frozen, please contact tenant of admin)
	MidMembershipHasBeenFrozenErr = tkErr.Error(MidMembershipHasBeenFrozenErrCode, MidMembershipHasBeenFrozenErrMsg)
	// 12000004(project (%s) not found)
	MidProjectNotFoundErr = tkErr.Error(MidProjectNotFoundErrCode, MidProjectNotFoundErrMsg)
	// 12000005(incorrect format of authentication)
	MidIncorrectFormatErr = tkErr.Error(MidIncorrectFormatErrCode, MidIncorrectFormatErrMsg)
	// 12000006(module category (%s) not found)
	MidModuleCategoryNotFoundErr = tkErr.Error(MidModuleCategoryNotFoundErrCode, MidModuleCategoryNotFoundErrCodeMsg)
	// 12000007(module (%s) not found)
	MidModuleNotFoundErr = tkErr.Error(MidModuleNotFoundErrCode, MidModuleNotFoundErrCodeMsg)
	// 12000008(module (%s) is read only for users of other projects)
	MidModuleIsReadOnlyErr = tkErr.Error(MidModuleIsReadOnlyErrCode, MidModuleIsReadOnlyErrMsg)
	// 12000009(application (%s) not found)
	MidApplicationNotFoundErr = tkErr.Error(MidApplicationNotFoundErrCode, MidApplicationNotFoundErrMsg)
	// 12000010(namespace not allow)
	MidNamespaceNotAllowErr = tkErr.Error(MidNamespaceNotAllowErrCode, MidNamespaceNotAllowErrMsg)
	// 12000011(application (%s) is read only for users of other projects)
	MidApplicationIsReadOnlyErr = tkErr.Error(MidApplicationIsReadOnlyErrCode, MidApplicationIsReadOnlyErrMsg)
	// 12000012(instance (%s) not found)
	MidInstanceNotFoundErr = tkErr.Error(MidInstanceNotFoundErrCode, MidInstanceNotFoundErrMsg)
	// 12000013(instance (%s) is read only for users of other projects)
	MidInstanceIsReadOnlyErr = tkErr.Error(MidInstanceIsReadOnlyErrCode, MidInstanceIsReadOnlyErrMsg)
)
