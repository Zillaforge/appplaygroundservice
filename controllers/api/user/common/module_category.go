package common

import (
	"context"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

type (
	ModuleCategory struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Creator     User   `json:"creator"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
		ModuleCount *int   `json:"moduleCount,omitempty"`
		_           struct{}
	}
)

func (data *ModuleCategory) ExtractByProto(ctx context.Context, input *pb.ModuleCategoryInfo) ModuleCategory {
	data.ID = input.ID
	data.Name = input.Name
	data.Description = input.Description

	data.Creator = User{
		ID: input.CreatorID,
	}
	data.Creator.Fill(ctx)

	data.CreatedAt = input.CreatedAt
	data.UpdatedAt = input.UpdatedAt

	return *data
}

func (data *ModuleCategory) ExtractByViewProto(ctx context.Context, input *pb.ModuleJoinModuleAclInfo) ModuleCategory {
	data.ID = input.ModuleCategoryID
	data.Name = input.ModuleCategoryName
	data.Description = input.ModuleCategoryDescription

	data.Creator = User{
		ID: input.ModuleCategoryCreatorID,
	}
	data.Creator.Fill(ctx)

	data.CreatedAt = input.ModuleCategoryCreatedAt
	data.UpdatedAt = input.ModuleCategoryUpdatedAt

	return *data
}
