package project

import (
	"AppPlaygroundService/storages/tables"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.ProjectCRUDControllerServer
type Method struct {
	// Embed UnsafeProjectCRUDControllerServer to have mustEmbedUnimplementedProjectCRUDControllerServer()
	pb.UnsafeProjectCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ProjectCRUDControllerServer = (*Method)(nil)

var empty = &emptypb.Empty{}

func (m Method) storage2pb(input *tables.Project) (output *pb.ProjectInfo) {
	return &pb.ProjectInfo{
		ID:        input.ID,
		CreatedAt: input.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: input.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
