package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1500xxxx: Storage

	StorageInternalServerErrCode              = 15000000
	StorageInternalServerErrMsg               = "internal server error"
	StorageOneOfResourcesNotFoundErrCode      = 15000001
	StorageOneOfResourcesNotFoundErrMsg       = "one of resource not found"
	StorageProjectExistErrCode                = 15000002
	StorageProjectExistErrMsg                 = "project exist"
	StorageProjectNotFoundErrCode             = 15000003
	StorageProjectNotFoundErrMsg              = "project not found"
	StorageProjectInUseErrCode                = 15000004
	StorageProjectInUseErrMsg                 = "project in use"
	StorageModuleCategoryExistErrCode         = 15000005
	StorageModuleCategoryExistErrMsg          = "module category exist"
	StorageModuleCategoryNotFoundErrCode      = 15000006
	StorageModuleCategoryNotFoundErrMsg       = "module category not found"
	StorageModuleCategoryInUseErrCode         = 15000007
	StorageModuleCategoryInUseErrMsg          = "module category in use"
	StorageModuleExistErrCode                 = 15000008
	StorageModuleExistErrMsg                  = "module exist"
	StorageModuleNotFoundErrCode              = 15000009
	StorageModuleNotFoundErrMsg               = "module not found"
	StorageModuleInUseErrCode                 = 15000010
	StorageModuleInUseErrMsg                  = "module in use"
	StorageModuleAclExistErrCode              = 15000011
	StorageModuleAclExistErrMsg               = "module-acl exist"
	StorageModuleAclNotFoundErrCode           = 15000012
	StorageModuleAclNotFoundErrMsg            = "module-acl not found"
	StorageModuleAclInUseErrCode              = 15000013
	StorageModuleAclInUseErrMsg               = "module-acl in use"
	StorageApplicationExistErrCode            = 15000014
	StorageApplicationExistErrMsg             = "application exist"
	StorageApplicationNotFoundErrCode         = 15000015
	StorageApplicationNotFoundErrMsg          = "application not found"
	StorageApplicationInUseErrCode            = 15000016
	StorageApplicationInUseErrMsg             = "application in use"
	StorageInstanceExistErrCode               = 15000017
	StorageInstanceExistErrMsg                = "instance exist"
	StorageInstanceNotFoundErrCode            = 15000018
	StorageInstanceNotFoundErrMsg             = "instance not found"
	StorageInstanceInUseErrCode               = 15000019
	StorageInstanceInUseErrMsg                = "instance in use"
	StorageModuleJoinModuleAclNotFoundErrCode = 15000020
	StorageModuleJoinModuleAclNotFoundErrMsg  = "module join module-acl not found"
	StorageMeteringExistErrCode               = 15000021
	StorageMeteringExistErrMsg                = "metering exist"
	StorageMeteringNotFoundErrCode            = 15000022
	StorageMeteringNotFoundErrMsg             = "metering not found"
	StorageAppCredentialExistErrCode          = 15000023
	StorageAppCredentialExistErrMsg           = "app credential exist"
	StorageAppCredentialInUseErrCode          = 15000024
	StorageAppCredentialInUseErrMsg           = "app credential in use"
	StorageAppCredentialNotFoundErrCode       = 15000025
	StorageAppCredentialNotFoundErrMsg        = "app credential not found"
)

var (
	// 1500xxxx: Storage

	// 15000000(internal server error)
	StorageInternalServerErr = tkErr.Error(StorageInternalServerErrCode, StorageInternalServerErrMsg)
	// 15000001 (one of resource not found)
	StorageOneOfResourcesNotFoundErr = tkErr.Error(StorageOneOfResourcesNotFoundErrCode, StorageOneOfResourcesNotFoundErrMsg)
	// 15000002 (project exist)
	StorageProjectExistErr = tkErr.Error(StorageProjectExistErrCode, StorageProjectExistErrMsg)
	// 15000003 (project not found)
	StorageProjectNotFoundErr = tkErr.Error(StorageProjectNotFoundErrCode, StorageProjectNotFoundErrMsg)
	// 15000004 (project in use)
	StorageProjectInUseErr = tkErr.Error(StorageProjectInUseErrCode, StorageProjectInUseErrMsg)
	// 15000005 (module category exist)
	StorageModuleCategoryExistErr = tkErr.Error(StorageModuleCategoryExistErrCode, StorageModuleCategoryExistErrMsg)
	// 15000006 (module category not found)
	StorageModuleCategoryNotFoundErr = tkErr.Error(StorageModuleCategoryNotFoundErrCode, StorageModuleCategoryNotFoundErrMsg)
	// 15000007 (module category in use)
	StorageModuleCategoryInUseErr = tkErr.Error(StorageModuleCategoryInUseErrCode, StorageModuleCategoryInUseErrMsg)
	// 15000008 (module exist)
	StorageModuleExistErr = tkErr.Error(StorageModuleExistErrCode, StorageModuleExistErrMsg)
	// 15000009 (module not found)
	StorageModuleNotFoundErr = tkErr.Error(StorageModuleNotFoundErrCode, StorageModuleNotFoundErrMsg)
	// 15000010 (module in use)
	StorageModuleInUseErr = tkErr.Error(StorageModuleInUseErrCode, StorageModuleInUseErrMsg)
	// 15000011 (module-acl exist)
	StorageModuleAclExistErr = tkErr.Error(StorageModuleAclExistErrCode, StorageModuleAclExistErrMsg)
	// 15000012 (module-acl not found)
	StorageModuleAclNotFoundErr = tkErr.Error(StorageModuleAclNotFoundErrCode, StorageModuleAclNotFoundErrMsg)
	// 15000013 (module-acl in use)
	StorageModuleAclInUseErr = tkErr.Error(StorageModuleAclInUseErrCode, StorageModuleAclInUseErrMsg)
	// 15000014 (application exist)
	StorageApplicationExistErr = tkErr.Error(StorageApplicationExistErrCode, StorageApplicationExistErrMsg)
	// 15000015 (application not found)
	StorageApplicationNotFoundErr = tkErr.Error(StorageApplicationNotFoundErrCode, StorageApplicationNotFoundErrMsg)
	// 15000016 (application in use)
	StorageApplicationInUseErr = tkErr.Error(StorageApplicationInUseErrCode, StorageApplicationInUseErrMsg)
	// 15000017 (instance exist)
	StorageInstanceExistErr = tkErr.Error(StorageInstanceExistErrCode, StorageInstanceExistErrMsg)
	// 15000018 (instance not found)
	StorageInstanceNotFoundErr = tkErr.Error(StorageInstanceNotFoundErrCode, StorageInstanceNotFoundErrMsg)
	// 15000019 (instance in use)
	StorageInstanceInUseErr = tkErr.Error(StorageInstanceInUseErrCode, StorageInstanceInUseErrMsg)
	// 15000020 (module join module-acl not found)
	StorageModuleJoinModuleAclNotFoundErr = tkErr.Error(StorageModuleJoinModuleAclNotFoundErrCode, StorageModuleJoinModuleAclNotFoundErrMsg)
	// 15000021(metering exist)
	StorageMeteringExistErr = tkErr.Error(StorageMeteringExistErrCode, StorageMeteringExistErrMsg)
	// 15000022(metering not found)
	StorageMeteringNotFoundErr = tkErr.Error(StorageMeteringNotFoundErrCode, StorageMeteringNotFoundErrMsg)
	// 15000021 (app credential exist)
	StorageAppCredentialExistErr = tkErr.Error(StorageAppCredentialExistErrCode, StorageAppCredentialExistErrMsg)
	// 15000022 (app credential in use)
	StorageAppCredentialInUseErr = tkErr.Error(StorageAppCredentialInUseErrCode, StorageAppCredentialInUseErrMsg)
	// 15000023 (app credential not found)
	StorageAppCredentialNotFoundErr = tkErr.Error(StorageAppCredentialNotFoundErrCode, StorageAppCredentialNotFoundErrMsg)
)
