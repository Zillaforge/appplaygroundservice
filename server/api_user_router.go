package server

import (
	cnt "AppPlaygroundService/constants"
	sysCtl "AppPlaygroundService/controllers/api/system"
	apiUserApplicationCtl "AppPlaygroundService/controllers/api/user/application"
	apiUserInstanceCtl "AppPlaygroundService/controllers/api/user/instance"
	apiUserModuleCtl "AppPlaygroundService/controllers/api/user/module"
	apiUserModuleCategoryCtl "AppPlaygroundService/controllers/api/user/module_category"
	mid "AppPlaygroundService/middlewares/api"

	"github.com/gin-gonic/gin"
	pbac "github.com/Zillaforge/toolkits/pbac/gin"
)

func enableUserAppPlaygroundServiceRouter(rg *gin.RouterGroup) {
	pbac.GET(rg, "version", sysCtl.GetPlainTextVersion, cnt.UserVersion.Name, false)

	rg.Use(mid.VerifyUserToken, mid.VerifyNamespace)

	projectID := rg.Group("project/:project-id", mid.VerifyMembership)
	{
		pbac.GET(projectID, "module-categories", apiUserModuleCategoryCtl.ListModuleCategories, cnt.UserListModuleCategories.Name, false)
		moduleCategoryID := projectID.Group("module-category/:module-category-id", mid.VerifyModuleCategory)
		{
			pbac.GET(moduleCategoryID, "", apiUserModuleCategoryCtl.GetModuleCategory, cnt.UserGetModuleCategory.Name, false)
			modules := moduleCategoryID.Group("modules", mid.SetHdrLanguage)
			{
				pbac.GET(modules, "", apiUserModuleCtl.ListModules, cnt.UserListModules.Name, false)
			}
		}

		moduleID := projectID.Group("module/:module-id", mid.SetHdrLanguage, mid.VerifyModuleHasInProject)
		{
			pbac.GET(moduleID, "", apiUserModuleCtl.GetModule, cnt.UserGetModule.Name, false)
		}

		applicationsRG := projectID.Group("applications", mid.SetHdrLanguage)
		{
			pbac.GET(applicationsRG, "", apiUserApplicationCtl.ListApplications, cnt.UserListApplications.Name, false)
		}

		applicationRG := projectID.Group("application", mid.SetHdrLanguage)
		{
			createApplicationRG := applicationRG.Group("", mid.ResourceReview)
			pbac.POST(createApplicationRG, "", apiUserApplicationCtl.CreateApplication, cnt.UserCreateApplication.Name, false)
			applicationIdRG := applicationRG.Group(":application-id", mid.VerifyApplicationHasInProject)
			{
				pbac.POST(applicationIdRG, "approve", apiUserApplicationCtl.ApproveApplication, cnt.UserApproveApplication.Name, false)
				pbac.POST(applicationIdRG, "reject", apiUserApplicationCtl.RejectApplication, cnt.UserRejectApplication.Name, false)

				applicationRoleCheckRG := applicationIdRG.Group("", mid.VerifyRolePermissionForApplication)
				pbac.GET(applicationRoleCheckRG, "", apiUserApplicationCtl.GetApplication, cnt.UserGetApplication.Name, false)
			}
		}

		applicationNoLangRG := projectID.Group("application/:application-id", mid.VerifyApplicationHasInProject, mid.VerifyRolePermissionForApplication)
		{
			pbac.DELETE(applicationNoLangRG, "", apiUserApplicationCtl.DeleteApplication, cnt.UserDeleteApplication.Name, false)
			pbac.GET(applicationNoLangRG, "logs", apiUserApplicationCtl.GetAppLogs, cnt.UserGetAppLogs.Name, false)
		}

		pbac.GET(projectID, "instances", apiUserInstanceCtl.ListInstances, cnt.UserListInstances.Name, false)
		instanceID := projectID.Group("instance/:instance-id", mid.VerifyInstanceHasInProject, mid.VerifyRolePermissionForInstance)
		{
			pbac.GET(instanceID, "", apiUserInstanceCtl.GetInstance, cnt.UserGetInstance.Name, false)
			pbac.PUT(instanceID, "", apiUserInstanceCtl.UpdateInstance, cnt.UserUpdateInstance.Name, false)

			floatingip := instanceID.Group("floatingip/associate")
			{
				pbac.POST(floatingip, "", apiUserInstanceCtl.AssociateFloatingIP, cnt.UserAssociateFloatingIP.Name, false)
				pbac.DELETE(floatingip, "", apiUserInstanceCtl.DisassociateFloatingIP, cnt.UserDisassociateFloatingIP.Name, false)
			}
		}
	}
}
