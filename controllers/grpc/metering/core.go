package metering

import (
	"AppPlaygroundService/storages/tables"
	"time"

	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.MeteringCRUDControllerServer
type Method struct {
	// Embed UnsafeMeteringCRUDControllerServer to have mustEmbedUnimplementedMeteringCRUDControllerServer()
	pb.UnsafeMeteringCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.MeteringCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.Metering) (output *pb.MeteringInfo) {
	output = &pb.MeteringInfo{
		ApplicationID: input.ApplicationID,
		Name:          input.Name,
		ProjectID:     input.ProjectID,
		Creator:       input.Creator,
		Instances:     input.Instances,
		CreatedAt:     input.CreatedAt.UTC().Format(time.RFC3339),
	}
	if input.EndedAt != nil {
		t := *input.EndedAt
		stringT := t.Format(time.RFC3339)
		output.EndedAt = &stringT
	}
	if input.LastPublishedAt != nil {
		t := *input.LastPublishedAt
		stringT := t.Format(time.RFC3339)
		output.LastPublishedAt = &stringT
	}
	return output
}
