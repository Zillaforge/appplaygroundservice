package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1302xxxx: User HTTP Controller

	UserAPIInternalServerErrCode           = 13020000
	UserAPIInternalServerErrMsg            = "internal server error"
	UserAPIQueryNotSupportErrCode          = 13020001
	UserAPIQueryNotSupportErrMsg           = "%s does not support %s"
	UserAPIIllegalWhereQueryFormatErrCode  = 13020002
	UserAPIIllegalWhereQueryFormatErrMsg   = "illegal query format with where"
	UserAPIApplicationExistErrCode         = 13020003
	UserAPIApplicationExistErrMsg          = "application (%s) exist"
	UserAPIInstanceNotFoundErrCode         = 13020004
	UserAPIInstanceNotFoundErrMsg          = "instance (%s) not found"
	UserAPIUnauthorizedOpErrCode           = 13020005
	UserAPIUnauthorizedOpErrMsg            = "unauthorized operation"
	UserAPIModuleNotFoundErrCode           = 13020006
	UserAPIModuleNotFoundErrMsg            = "module (%s) is not exists"
	UserAPIGetModuleQuestionsFailedErrCode = 13020007
	UserAPIGetModuleQuestionsFailedErrMsg  = "get module questions failed"
	UserAPIApplicationIsProcessingErrCode  = 13020008
	UserAPIApplicationIsProcessingErrMsg   = "application is processing"
	UserAPIProjectNotFoundErrCode          = 13020009
	UserAPIProjectNotFoundErrMsg           = "project (%s) not found"
	UserAPIInternalServerErrWithInnerCode  = 13020010
	UserAPIInternalServerErrWithInnerMsg   = "internal server error: [%s]"
	UserAPIQuizModuleErrCode               = 13020011
	UserAPIQuizModuleErrMsg                = "%s"
	UserAPIFloatingIPNotFoundErrCode       = 13020012
	UserAPIFloatingIPNotFoundErrMsg        = "floating ip (%s) not found"
	UserAPIInstanceAlreadyHasFIPErrCode    = 13020013
	UserAPIInstanceAlreadyHasFIErrMsg      = "instance already has a floating ip"
	UserAPIFIPCannotBeUseErrCode           = 13020014
	UserAPIFIPCannotBeUseErrMsg            = "floating ip can not be use"
	UserAPIInstanceHasNoFIPErrCode         = 13020015
	UserAPIInstanceHasNoFIPErrMsg          = "instance has no floating ip"
)

var (
	// 1302xxxx: User HTTP Controller

	// 13020000(internal server error)
	UserAPIInternalServerErr = tkErr.Error(UserAPIInternalServerErrCode, UserAPIInternalServerErrMsg)
	// 13020001(%s does not support %s)
	UserAPIQueryNotSupportErr = tkErr.Error(UserAPIQueryNotSupportErrCode, UserAPIQueryNotSupportErrMsg)
	// 13020002(illegal query format with where)
	UserAPIIllegalWhereQueryFormatErr = tkErr.Error(UserAPIIllegalWhereQueryFormatErrCode, UserAPIIllegalWhereQueryFormatErrMsg)
	// 13020003(application (%s) exist)
	UserAPIApplicationExistErr = tkErr.Error(UserAPIApplicationExistErrCode, UserAPIApplicationExistErrMsg)
	// 13020004(instance (%s) not found)
	UserAPIInstanceNotFoundErr = tkErr.Error(UserAPIInstanceNotFoundErrCode, UserAPIInstanceNotFoundErrMsg)
	// 13020005(unauthorized operation)
	UserAPIUnauthorizedOpErr = tkErr.Error(UserAPIUnauthorizedOpErrCode, UserAPIUnauthorizedOpErrMsg)
	// 13020006(module (%d) is not exists)
	UserAPIModuleNotFoundErr = tkErr.Error(UserAPIModuleNotFoundErrCode, UserAPIModuleNotFoundErrMsg)
	// 13020007(get module questions failed)
	UserAPIGetModuleQuestionsFailedErr = tkErr.Error(UserAPIGetModuleQuestionsFailedErrCode, UserAPIGetModuleQuestionsFailedErrMsg)
	// 13020008(application is processing)
	UserAPIApplicationIsProcessingErr = tkErr.Error(UserAPIApplicationIsProcessingErrCode, UserAPIApplicationIsProcessingErrMsg)
	// 13010009(project not found)
	UserAPIProjectNotFoundErr = tkErr.Error(UserAPIProjectNotFoundErrCode, UserAPIProjectNotFoundErrMsg)
	// 13020010(internal server error: [%s])
	UserAPIInternalServerErrWithInner = tkErr.Error(UserAPIInternalServerErrWithInnerCode, UserAPIInternalServerErrWithInnerMsg)
	// 13020011(%s) ex. quiz module error: [flavor (3fe78aca-3ac7-4051-a1f0-5baf3d20443f) not found]
	UserAPIQuizModuleErr = tkErr.Error(UserAPIQuizModuleErrCode, UserAPIQuizModuleErrMsg)
	// 13020012(floating ip (%s) not found)
	UserAPIFloatingIPNotFoundErr = tkErr.Error(UserAPIFloatingIPNotFoundErrCode, UserAPIFloatingIPNotFoundErrMsg)
	// 13020013(instance already has a floating ip)
	UserAPIInstanceAlreadyHasFIPErr = tkErr.Error(UserAPIInstanceAlreadyHasFIPErrCode, UserAPIInstanceAlreadyHasFIErrMsg)
	// 13020014(floating ip can not be use)
	UserAPIFIPCannotBeUseErr = tkErr.Error(UserAPIFIPCannotBeUseErrCode, UserAPIFIPCannotBeUseErrMsg)
	// 13020015(instance has no floating ip)
	UserAPIInstanceHasNoFIPErr = tkErr.Error(UserAPIInstanceHasNoFIPErrCode, UserAPIInstanceHasNoFIPErrMsg)
)
