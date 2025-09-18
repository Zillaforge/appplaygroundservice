package server

import (
	cnt "AppPlaygroundService/constants"
	apiAdminApplicationCtl "AppPlaygroundService/controllers/api/admin/application"
	apiAdminInstanceCtl "AppPlaygroundService/controllers/api/admin/instance"
	apiAdminModuleCtl "AppPlaygroundService/controllers/api/admin/module"
	apiAdminModuleACLCtl "AppPlaygroundService/controllers/api/admin/module_acl"
	apiAdminModuleCategoryCtl "AppPlaygroundService/controllers/api/admin/module_category"
	sysCtl "AppPlaygroundService/controllers/api/system"
	mid "AppPlaygroundService/middlewares/api"
	"fmt"

	"github.com/gin-gonic/gin"
	pbac "github.com/Zillaforge/toolkits/pbac/gin"
)

func enableAdminAppPlaygroundServiceRouter(rg *gin.RouterGroup) {
	rg.Use(mid.VerifyAdminAuthentication)
	system := rg.Group("system")
	{
		pbac.GET(system, "versions", sysCtl.GetDetailVersions, cnt.AdminGetDetailVersions.Name, true)
		// api/v1/admin/system/configurations
		pbac.GET(system, "configurations", sysCtl.GetSystemConfigurations, cnt.AdminGetSystemConfigurations.Name, true)
		logs := system.Group("logs")
		{
			pbac.GET(logs, "", sysCtl.ListLogs, cnt.AdminListLogs.Name, true)
			pbac.GET(logs, "download", sysCtl.DownloadLog, cnt.AdminDownloadLog.Name, true)
		}
	}

	// application
	applicationRGs := rg.Group("applications", mid.SetHdrLanguage, mid.VerifyNamespace)
	{
		pbac.GET(applicationRGs, "", apiAdminApplicationCtl.ListApplications, cnt.AdminListApplications.Name, true)
	}
	applicationRG := rg.Group("application", mid.SetHdrLanguage)
	{
		pbac.POST(applicationRG, "", apiAdminApplicationCtl.CreateApplication, cnt.AdminCreateApplication.Name, true)
		applicationIdRG := applicationRG.Group(fmt.Sprintf(":%s", cnt.ParamApplicationID), mid.VerifyApplication)
		{
			pbac.GET(applicationIdRG, "", apiAdminApplicationCtl.GetApplication, cnt.AdminGetApplication.Name, true)

		}
	}
	applicationNoLangRG := rg.Group(fmt.Sprintf("application/:%s", cnt.ParamApplicationID), mid.VerifyApplication)
	{
		pbac.DELETE(applicationNoLangRG, "", apiAdminApplicationCtl.DeleteApplication, cnt.AdminDeleteApplication.Name, true)
		pbac.GET(applicationNoLangRG, "logs", apiAdminApplicationCtl.GetAppLogs, cnt.AdminGetAppLogs.Name, true)
	}

	// module-category
	moduleCategoryRGs := rg.Group("module-categories")
	{
		pbac.GET(moduleCategoryRGs, "", apiAdminModuleCategoryCtl.ListModuleCategories, cnt.AdminListModuleCategories.Name, true)
	}
	moduleCategoryRG := rg.Group("module-category")
	{
		pbac.POST(moduleCategoryRG, "", apiAdminModuleCategoryCtl.CreateModuleCategory, cnt.AdminCreateModuleCategory.Name, true)
		moduleCategoryIdRG := moduleCategoryRG.Group(fmt.Sprintf(":%s", cnt.ParamModuleCategoryID), mid.VerifyModuleCategory)
		{
			pbac.GET(moduleCategoryIdRG, "", apiAdminModuleCategoryCtl.GetModuleCategory, cnt.AdminGetModuleCategory.Name, true)
			pbac.DELETE(moduleCategoryIdRG, "", apiAdminModuleCategoryCtl.DeleteModuleCategory, cnt.AdminDeleteModuleCategory.Name, true)
			pbac.PUT(moduleCategoryIdRG, "", apiAdminModuleCategoryCtl.UpdateModuleCategory, cnt.AdminDeleteModuleCategory.Name, true)
		}
	}

	// module
	moduleCRGs := rg.Group("modules", mid.SetHdrLanguage)
	{
		pbac.GET(moduleCRGs, "", apiAdminModuleCtl.ListModules, cnt.AdminListModules.Name, true)
	}
	moduleRG := rg.Group("module", mid.SetHdrLanguage)
	{
		pbac.POST(moduleRG, "", apiAdminModuleCtl.CreateModule, cnt.AdminCreateModule.Name, true)
		moduleIdRG := moduleRG.Group(fmt.Sprintf(":%s", cnt.ParamModuleID), mid.VerifyModule)
		{
			pbac.GET(moduleIdRG, "", apiAdminModuleCtl.GetModule, cnt.AdminGetModule.Name, true)
			pbac.PUT(moduleIdRG, "", apiAdminModuleCtl.UpdateModule, cnt.AdminCreateModule.Name, true)
		}
	}
	moduleNoLangRG := rg.Group(fmt.Sprintf("module/:%s", cnt.ParamModuleID), mid.VerifyModule)
	{
		pbac.DELETE(moduleNoLangRG, "", apiAdminModuleCtl.DeleteModule, cnt.AdminDeleteModule.Name, true)

		pbac.GET(moduleNoLangRG, "/acl", apiAdminModuleACLCtl.GetModuleACL, cnt.AdminGetModuleACL.Name, true)
		pbac.PUT(moduleNoLangRG, "/acl", apiAdminModuleACLCtl.UpdateModuleACL, cnt.AdminUpdateModuleACL.Name, true)
	}

	// instance
	instanceCRGs := rg.Group("instances")
	{
		pbac.GET(instanceCRGs, "", apiAdminInstanceCtl.ListInstances, cnt.AdminListInstances.Name, true)
	}
	instanceRG := rg.Group("instance")
	{
		instanceIdRG := instanceRG.Group(fmt.Sprintf(":%s", cnt.ParamInstanceID), mid.VerifyInstance)
		{
			pbac.GET(instanceIdRG, "", apiAdminInstanceCtl.GetInstance, cnt.AdminGetInstance.Name, true)
			pbac.PUT(instanceIdRG, "", apiAdminInstanceCtl.UpdateInstance, cnt.AdminUpdateInstance.Name, true)

			floatingipRG := instanceIdRG.Group("floatingip/associate")
			{
				pbac.POST(floatingipRG, "", apiAdminInstanceCtl.AssociateFloatingIP, cnt.AdminAssociateFloatingIP.Name, true)
				pbac.DELETE(floatingipRG, "", apiAdminInstanceCtl.DisassociateFloatingIP, cnt.AdminDisassociateFloatingIP.Name, true)
			}
		}
	}
}
