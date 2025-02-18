package constants

import (
	"fmt"
	"strconv"
)

/*
ServiceID: an int [ 0 - 99 ]
Service.ID 註冊後就無法更改，要是要更改請要通知 UI/UX Team或 Nokia Du(杜岳霖)
註冊 ServiceID請至 gitlab wiki 表格
(<gitlab-url>/aes/pegasus-guide/pegasuscontainerimagebuilder/-/wikis/develop-environment-(docker-compose)#serviceid-table)

actionIDRange: 代表 ActionID 的數量上限
如果 actionIDRange 是 1000, 代表 ActionID: [ 0 - 999 ] + serviceID * 1000
Ex: serviceID = 14, actionIDRange = 1000, ActionID: [ 14000 - 14999 ]
*/
const (
	serviceID     = 27
	actionIDRange = 1000
)

var (
	accessLoggerInfoNameMap map[string]*AccessLoggerInfo = make(map[string]*AccessLoggerInfo)
	accessLoggerInfoIDMap   map[int]*AccessLoggerInfo    = make(map[int]*AccessLoggerInfo)
)

type AccessLoggerInfo struct {
	// ID 固定正數值且不可以重複
	ID int
	// Name 須具備標準格式
	Name string
}

func new(id int, name string) AccessLoggerInfo {
	if id >= actionIDRange || id < 0 {
		id = id % actionIDRange
	}
	id = id + serviceID*actionIDRange
	alInfo := AccessLoggerInfo{id, name}

	if _, exist := accessLoggerInfoNameMap[alInfo.Name]; exist {
		panic(fmt.Sprintf("access logger name is duplicate: %s", alInfo.Name))
	}
	if _, exist := accessLoggerInfoIDMap[alInfo.ID]; exist {
		panic(fmt.Sprintf("access logger id is duplicate: %d", alInfo.ID))
	}
	accessLoggerInfoNameMap[alInfo.Name] = &alInfo
	accessLoggerInfoIDMap[alInfo.ID] = &alInfo
	return alInfo
}

func GetAccessLoggerInfo(name string) *AccessLoggerInfo {
	if v, ok := accessLoggerInfoNameMap[name]; ok {
		return v
	}
	return nil
}

func InsertNewAccessLoggerInfo(actionName string, actionID int) AccessLoggerInfo {
	return new(actionID, actionName)
}

func GetAccessLoggerServiceIDStr() string {
	return strconv.Itoa(serviceID)
}

// AccessLoggerInfo 是用來收集 RESTful API的事件名稱與數字編號的關係。
//
// AccessLoggerInfo.Name須有特定格式來作為不同資源區分。格式會採用如下
//
// <category name>:<operation or action name>
//
// 目前已知的 <category name> 有： user, admin, system分別表示 使用者RESTful
// 管理者RESTful以及系統類型
//
// 特別注意：
// 1). AccessLoggerInfo.ID 註冊後就無法更改，要是要更改請要通知 UI/UX Team或 Nokia Du(杜岳霖)
//

var (
	// Category name : system
	// CategoryName=system通常運作於系統本身或使用者/管理者的間接行為。
	// e.g. SysAsyncPublishProcess AccessLoggerInfo = new(withPrefix(999), "system:SysAsyncPublishProcess")

	// Category name: user
	// 運作於使用者 RESTful API的動作，都需要為 `user:`開頭的前綴
	// e.g. UserVersion AccessLoggerInfo = new(withPrefix(0), "user:Version")
	UserListModuleCategories   AccessLoggerInfo = new(0, "user:ListModuleCategories")
	UserGetModuleCategory      AccessLoggerInfo = new(1, "user:GetModuleCategory")
	UserListModules            AccessLoggerInfo = new(2, "user:ListModules")
	UserGetModule              AccessLoggerInfo = new(3, "user:GetModule")
	UserListApplications       AccessLoggerInfo = new(4, "user:ListApplications")
	UserGetApplication         AccessLoggerInfo = new(5, "user:GetApplication")
	UserCreateApplication      AccessLoggerInfo = new(6, "user:CreateApplication")
	UserDeleteApplication      AccessLoggerInfo = new(7, "user:DeleteApplication")
	UserApproveApplication     AccessLoggerInfo = new(8, "user:ApproveApplication")
	UserListInstances          AccessLoggerInfo = new(9, "user:ListInstances")
	UserGetInstance            AccessLoggerInfo = new(10, "user:GetInstance")
	UserUpdateInstance         AccessLoggerInfo = new(11, "user:UpdateInstance")
	UserRejectApplication      AccessLoggerInfo = new(12, "user:RejectApplication")
	UserAssociateFloatingIP    AccessLoggerInfo = new(13, "user:AssociateFloatingIP")
	UserDisassociateFloatingIP AccessLoggerInfo = new(14, "user:DissassociateFloatingIP")
	UserGetAppLogs             AccessLoggerInfo = new(15, "user:GetAppLogs")
	UserVersion                AccessLoggerInfo = new(16, "user:Version")

	// Category name: admin
	// 運作於管理者 RESTful API的動作，都需要為 `admin:`開頭的前綴
	// Start ID from 300
	AdminGetDetailVersions       AccessLoggerInfo = new(300, "admin:GetDetailVersions")
	AdminGetSystemConfigurations AccessLoggerInfo = new(301, "admin:GetSystemConfigurations")
	AdminListLogs                AccessLoggerInfo = new(302, "admin:ListLogs")
	AdminDownloadLog             AccessLoggerInfo = new(303, "admin:DownloadLog")
	AdminListApplications        AccessLoggerInfo = new(304, "admin:ListApplications")
	AdminGetApplication          AccessLoggerInfo = new(305, "admin:GetApplication")
	AdminDeleteApplication       AccessLoggerInfo = new(306, "admin:DeleteApplication")
	AdminCreateApplication       AccessLoggerInfo = new(307, "admin:CreateApplication")
	AdminListModuleCategories    AccessLoggerInfo = new(308, "admin:ListModuleCategories")
	AdminCreateModuleCategory    AccessLoggerInfo = new(309, "admin:CreateModuleCategory")
	AdminGetModuleCategory       AccessLoggerInfo = new(310, "admin:GetModuleCategory")
	AdminDeleteModuleCategory    AccessLoggerInfo = new(311, "admin:DeleteModuleCategory")
	AdminUpdateModuleCategory    AccessLoggerInfo = new(312, "admin:UpdateModuleCategory")
	AdminListModules             AccessLoggerInfo = new(313, "admin:ListModules")
	AdminGetModule               AccessLoggerInfo = new(314, "admin:GetModule")
	AdminCreateModule            AccessLoggerInfo = new(315, "admin:CreateModule")
	AdminUpdateModule            AccessLoggerInfo = new(316, "admin:UpdateModule")
	AdminDeleteModule            AccessLoggerInfo = new(317, "admin:DeleteModule")
	AdminGetModuleACL            AccessLoggerInfo = new(318, "admin:GetModuleACL")
	AdminUpdateModuleACL         AccessLoggerInfo = new(319, "admin:UpdateModuleACL")
	AdminListInstances           AccessLoggerInfo = new(320, "admin:ListInstances")
	AdminGetInstance             AccessLoggerInfo = new(321, "admin:GetInstance")
	AdminUpdateInstance          AccessLoggerInfo = new(322, "admin:UpdateInstance")
	AdminAssociateFloatingIP     AccessLoggerInfo = new(323, "admin:AssociateFloatingIP")
	AdminDisassociateFloatingIP  AccessLoggerInfo = new(324, "admin:DissassociateFloatingIP")
	AdminGetAppLogs              AccessLoggerInfo = new(325, "admin:GetAppLogs")
)
