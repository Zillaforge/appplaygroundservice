package common

import (
	auth "AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	"context"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

type ModuleCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     User   `json:"creator"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ModuleCount *int   `json:"moduleCount,omitempty"`
	_           struct{}
}

type ModuleCategoryInterface interface {
	ExtractByProto(ctx context.Context, proto *pb.ModuleCategoryInfo)
	ExtractByViewProto(ctx context.Context, proto *pb.ModuleJoinModuleAclInfo)
}

func (a *ModuleCategory) ExtractByProto(ctx context.Context, proto *pb.ModuleCategoryInfo) {
	if proto == nil {
		return
	}
	a.ID = proto.ID
	a.Name = proto.Name
	a.Description = proto.Description

	getUserInput := &authCom.GetUserInput{
		ID:        proto.CreatorID,
		Cacheable: true,
	}
	getCreatorOutput, getCreatorErr := auth.Use().GetUser(ctx, getUserInput)
	if getCreatorErr == nil {
		a.Creator = User{
			ID:      getCreatorOutput.ID,
			Account: getCreatorOutput.Account,
			Name:    getCreatorOutput.DisplayName,
		}
	}

	a.CreatedAt = proto.CreatedAt
	a.UpdatedAt = proto.UpdatedAt
}

func (a *ModuleCategory) ExtractByViewProto(ctx context.Context, proto *pb.ModuleJoinModuleAclInfo) {
	if proto == nil {
		return
	}
	a.ID = proto.ModuleCategoryID
	a.Name = proto.ModuleCategoryName
	a.Description = proto.ModuleCategoryDescription

	getUserInput := &authCom.GetUserInput{
		ID:        proto.ModuleCategoryCreatorID,
		Cacheable: true,
	}
	getCreatorOutput, getCreatorErr := auth.Use().GetUser(ctx, getUserInput)
	if getCreatorErr == nil {
		a.Creator = User{
			ID:      getCreatorOutput.ID,
			Account: getCreatorOutput.Account,
			Name:    getCreatorOutput.DisplayName,
		}
	}

	a.CreatedAt = proto.ModuleCategoryCreatedAt
	a.UpdatedAt = proto.ModuleCategoryUpdatedAt
}
