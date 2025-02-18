package app_credential

import (
	"AppPlaygroundService/storages/tables"
	"time"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.AppCredentialCRUDControllerServer
type Method struct {
	// Embed UnsafeAppCredentialCRUDControllerServer to have mustEmbedUnimplementedAppCredentialCRUDControllerServer()
	pb.UnsafeAppCredentialCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.AppCredentialCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.AppCredential) (output *pb.AppCredentialInfo) {
	output = &pb.AppCredentialInfo{
		ID:        input.ID,
		Name:      input.Name,
		UserID:    input.UserID,
		ProjectID: input.ProjectID,
		Secret:    input.Secret,
		Namespace: input.Namespace,
		CreatedAt: input.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: input.UpdatedAt.UTC().Format(time.RFC3339),
	}
	return
}
