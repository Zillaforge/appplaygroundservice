package common

import (
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

type ModuleACL struct {
	ModuleID   string   `json:"moduleId"`
	ProjectIDs []string `json:"projectIds"`
	_          struct{}
}

type ModuleACLInterface interface {
	ExtractByProto(proto []*pb.ModuleAclInfo)
}

func (a *ModuleACL) ExtractByProto(proto []*pb.ModuleAclInfo) {
	if proto == nil {
		return
	}

	a.ModuleID = proto[0].ModuleID
	for _, moduleAcl := range proto {
		if moduleAcl.ProjectID == "" && len(proto) == 1 {
			a.ProjectIDs = []string{}
			return
		}
		a.ProjectIDs = append(a.ProjectIDs, moduleAcl.ProjectID)
	}
}
