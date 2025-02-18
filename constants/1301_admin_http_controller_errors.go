package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1301xxxx: Admin HTTP Controller

	AdminControllerLogFileNotFoundErrCode      = 13010000
	AdminControllerLogFileNotFoundErrMsg       = "log file not found"
	AdminControllerSenderMalformedInputErrCode = 13010001
	AdminControllerSenderMalformedInputErrMsg  = "input format is invalid (%s)"
	AdminAPIApplicationAlreadyExistErrCode     = 13010002
	AdminAPIApplicationAlreadyExistErrMsg      = "application already exist"
	AdminAPIApplicationNotFoundErrCode         = 13010003
	AdminAPIApplicationNotFoundErrMsg          = "application not found"
	AdminAPIModuleCategoryAlreadyExistErrCode  = 13010004
	AdminAPIModuleCategoryAlreadyExistErrMsg   = "module category already exist"
	AdminAPIModuleCategoryNotFoundErrCode      = 13010005
	AdminAPIModuleCategoryNotFoundErrMsg       = "module category not found"
	AdminAPIModuleCategoryAlreadyInUseErrCode  = 13010006
	AdminAPIModuleCategoryAlreadyInUseErrMsg   = "module category already in use"
	AdminAPIModuleNotFoundErrCode              = 13010007
	AdminAPIModuleNotFoundErrMsg               = "module (%s) not found"
	AdminAPIModuleAlreadyExistErrCode          = 13010008
	AdminAPIModuleAlreadyExistErrMsg           = "module already exist"
	AdminAPIModuleAlreadyInUseErrCode          = 13010009
	AdminAPIModuleAlreadyInUseErrMsg           = "module already in use"
	AdminAPIInternalServerErrCode              = 13010010
	AdminAPIInternalServerErrMsg               = "internal server error"
	AdminAPIModuleACLNotFoundErrCode           = 13010011
	AdminAPIModuleACLNotFoundErrMsg            = "module acl not found"
	AdminAPIProjectNotFoundErrCode             = 13010012
	AdminAPIProjectNotFoundErrMsg              = "project (%s) not found"
	AdminAPIInstanceNotFoundErrCode            = 13010013
	AdminAPIInstanceNotFoundErrMsg             = "instance not found"
	AdminAPIUserNotFoundErrCode                = 13010014
	AdminAPIUserNotFoundErrMsg                 = "user (%s) not found"
	AdminAPIGetModuleQuestionsFailedErrCode    = 13010015
	AdminAPIGetModuleQuestionsFailedErrMsg     = "get module questions failed"
	AdminAPIApplicationIsProcessingErrCode     = 13010016
	AdminAPIApplicationIsProcessingErrMsg      = "application is processing"
	AdminAPIInternalServerErrWithInnerCode     = 13010017
	AdminAPIInternalServerErrWithInnerMsg      = "internal server error: [%s]"
	AdminAPIQuizModuleErrCode                  = 13010018
	AdminAPIQuizModuleErrMsg                   = "%s"
	AdminAPIFloatingIPNotFoundErrCode          = 13010019
	AdminAPIFloatingIPNotFoundErrMsg           = "floating ip (%s) not found"
	AdminAPIInstanceAlreadyHasFIPErrCode       = 13010020
	AdminAPIInstanceAlreadyHasFIErrMsg         = "instance already has a floating ip"
	AdminAPIFIPCannotBeUseErrCode              = 13010021
	AdminAPIFIPCannotBeUseErrMsg               = "floating ip can not be use"
	AdminAPIInstanceHasNoFIPErrCode            = 13010022
	AdminAPIInstanceHasNoFIEMsg                = "instance has no floating ip"
)

var (
	// 1301xxxx: Admin HTTP Controller

	// 13010000(log file not found)
	AdminControllerLogFileNotFoundErr = tkErr.Error(AdminControllerLogFileNotFoundErrCode, AdminControllerLogFileNotFoundErrMsg)
	// 13010001(input format is invalid)
	AdminControllerSenderMalformedInputErr = tkErr.Error(AdminControllerSenderMalformedInputErrCode, AdminControllerSenderMalformedInputErrMsg)
	// 13010002(application already exist)
	AdminAPIApplicationAlreadyExistErr = tkErr.Error(AdminAPIApplicationAlreadyExistErrCode, AdminAPIApplicationAlreadyExistErrMsg)
	// 13010003(application not found)
	AdminAPIApplicationNotFoundErr = tkErr.Error(AdminAPIApplicationNotFoundErrCode, AdminAPIApplicationNotFoundErrMsg)
	// 13010004(module category already exist)
	AdminAPIModuleCategoryAlreadyExistErr = tkErr.Error(AdminAPIModuleCategoryAlreadyExistErrCode, AdminAPIModuleCategoryAlreadyExistErrMsg)
	// 13010005(module category not found)
	AdminAPIModuleCategoryNotFoundErr = tkErr.Error(AdminAPIModuleCategoryNotFoundErrCode, AdminAPIModuleCategoryNotFoundErrMsg)
	// 13010006(module category already in use)
	AdminAPIModuleCategoryAlreadyInUseErr = tkErr.Error(AdminAPIModuleCategoryAlreadyInUseErrCode, AdminAPIModuleCategoryAlreadyInUseErrMsg)
	// 13010007(module not found)
	AdminAPIModuleNotFoundErr = tkErr.Error(AdminAPIModuleNotFoundErrCode, AdminAPIModuleNotFoundErrMsg)
	// 13010008(module already exist)
	AdminAPIModuleAlreadyExistErr = tkErr.Error(AdminAPIModuleAlreadyExistErrCode, AdminAPIModuleAlreadyExistErrMsg)
	// 13010009(module already in use)
	AdminAPIModuleAlreadyInUseErr = tkErr.Error(AdminAPIModuleAlreadyInUseErrCode, AdminAPIModuleAlreadyInUseErrMsg)
	// 13010010(internal server error)
	AdminAPIInternalServerErr = tkErr.Error(AdminAPIInternalServerErrCode, AdminAPIInternalServerErrMsg)
	// 13010011(module acl not found)
	AdminAPIModuleACLNotFoundErr = tkErr.Error(AdminAPIModuleACLNotFoundErrCode, AdminAPIModuleACLNotFoundErrMsg)
	// 13010012(project not found)
	AdminAPIProjectNotFoundErr = tkErr.Error(AdminAPIProjectNotFoundErrCode, AdminAPIProjectNotFoundErrMsg)
	// 13010013(instance not found)
	AdminAPIInstanceNotFoundErr = tkErr.Error(AdminAPIInstanceNotFoundErrCode, AdminAPIInstanceNotFoundErrMsg)
	// 13010014(user not found)
	AdminAPIUserNotFoundErr = tkErr.Error(AdminAPIUserNotFoundErrCode, AdminAPIUserNotFoundErrMsg)
	// 13010015(get module questions failed)
	AdminAPIGetModuleQuestionsFailedErr = tkErr.Error(AdminAPIGetModuleQuestionsFailedErrCode, AdminAPIGetModuleQuestionsFailedErrMsg)
	// 13010016(application is processing)
	AdminAPIApplicationIsProcessingErr = tkErr.Error(AdminAPIApplicationIsProcessingErrCode, AdminAPIApplicationIsProcessingErrMsg)
	// 13010017(internal server error: [%s])
	AdminAPIInternalServerErrWithInner = tkErr.Error(AdminAPIInternalServerErrWithInnerCode, AdminAPIInternalServerErrWithInnerMsg)
	// 13010018(%s) ex. quiz module error: [flavor (3fe78aca-3ac7-4051-a1f0-5baf3d20443f) not found]
	AdminAPIQuizModuleErr = tkErr.Error(AdminAPIQuizModuleErrCode, AdminAPIQuizModuleErrMsg)
	// 13010019(floating ip not found)
	AdminAPIFloatingIPNotFoundErr = tkErr.Error(AdminAPIFloatingIPNotFoundErrCode, AdminAPIFloatingIPNotFoundErrMsg)
	// 13010020(instance already has a floating ip)
	AdminAPIInstanceAlreadyHasFIPErr = tkErr.Error(AdminAPIInstanceAlreadyHasFIPErrCode, AdminAPIInstanceAlreadyHasFIErrMsg)
	// 13010021(floating ip can not be use)
	AdminAPIFIPCannotBeUseErr = tkErr.Error(AdminAPIFIPCannotBeUseErrCode, AdminAPIFIPCannotBeUseErrMsg)
	// 13010022(instance has no floating ip)
	AdminAPIInstanceHasNoFIPErr = tkErr.Error(AdminAPIInstanceHasNoFIPErrCode, AdminAPIInstanceHasNoFIEMsg)
)
