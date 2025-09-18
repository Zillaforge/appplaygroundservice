package modulecategory

import (
	"AppPlaygroundService/storages/tables"
	"time"

	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.ModuleCategoryCRUDControllerServer
type Method struct {
	// Embed UnsafeModuleCategoryCRUDControllerServer to have mustEmbedUnimplementedModuleCategoryCRUDControllerServer()
	pb.UnsafeModuleCategoryCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ModuleCategoryCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.ModuleCategory) (output *pb.ModuleCategoryInfo) {
	return &pb.ModuleCategoryInfo{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		CreatorID:   input.CreatorID,
		CreatedAt:   input.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   input.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
