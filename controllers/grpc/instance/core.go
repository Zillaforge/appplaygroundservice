package instance

import (
	"AppPlaygroundService/storages/tables"
	"time"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.InstanceCRUDControllerServer
type Method struct {
	// Embed UnsafeInstanceCRUDControllerServer to have mustEmbedUnimplementedInstanceCRUDControllerServer()
	pb.UnsafeInstanceCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.InstanceCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.Instance) (output *pb.InstanceDetail) {
	output = &pb.InstanceDetail{
		Instance: &pb.InstanceInfo{
			ID:                input.ID,
			Name:              input.Name,
			ApplicationID:     input.ApplicationID,
			ProjectID:         input.ProjectID,
			ReferenceID:       input.ReferenceID,
			Extra:             input.Extra,
			FloatingIPID:      input.FloatingIPID,
			FloatingIPAddress: input.FloatingIPAddress,
			CreatedAt:         input.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:         input.UpdatedAt.UTC().Format(time.RFC3339),
		},
		Application: &pb.ApplicationInfo{
			ID:          input.Application.ID,
			Name:        input.Application.Name,
			Description: input.Application.Description,
			ModuleID:    input.Application.ModuleID,
			State:       input.Application.State,
			Answers:     input.Application.Answers,
			Namespace:   input.Application.Namespace,
			Shiftable:   input.Application.Shiftable,
			ProjectID:   input.Application.ProjectID,
			CreatorID:   input.Application.CreatorID,
			UpdaterID:   input.Application.UpdaterID,
			CreatedAt:   input.Application.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   input.Application.UpdatedAt.UTC().Format(time.RFC3339),
		},
	}

	return
}
