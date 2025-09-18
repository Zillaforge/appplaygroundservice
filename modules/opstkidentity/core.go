package opstkidentity

import (
	"github.com/Zillaforge/toolkits/mviper"
)

var instance *OpstkIdentity

type OpstkIdentity struct {
	PidSource string
	UidSource string
}

func Init() *OpstkIdentity {
	instance := &OpstkIdentity{
		PidSource: mviper.GetString("app_playground_service.opstkidentity.pid_source"),
		UidSource: mviper.GetString("app_playground_service.opstkidentity.uid_source"),
	}
	return instance
}

func Use() *OpstkIdentity {
	if instance == nil {
		instance = Init()
	}
	return instance
}
