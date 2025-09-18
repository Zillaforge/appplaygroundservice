package moduleacl

import (
	"AppPlaygroundService/storages/tables"

	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.ModuleAclCRUDControllerServer
type Method struct {
	// Embed UnsafeModuleAclCRUDControllerServer to have mustEmbedUnimplementedModuleAclCRUDControllerServer()
	pb.UnsafeModuleAclCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ModuleAclCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.ModuleAcl) (output *pb.ModuleAclInfo) {
	return &pb.ModuleAclInfo{
		ID:        input.ID,
		ModuleID:  input.ModuleID,
		ProjectID: input.ProjectID,
	}
}
