package common

import (
	"context"
	"encoding/json"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

type (
	Module struct {
		ID             string                 `json:"id"`
		Name           string                 `json:"name"`
		Description    string                 `json:"description"`
		State          string                 `json:"state"`
		Questions      map[string]interface{} `json:"questions"`
		ModuleCategory Resource               `json:"moduleCategory"`
		Creator        User                   `json:"creator"`
		CreatedAt      string                 `json:"createdAt"`
		UpdatedAt      string                 `json:"updatedAt"`
		_              struct{}
	}
)

func (data *Module) ExtractByProto(ctx context.Context, input *pb.ModuleDetail) Module {
	if input == nil {
		return Module{}
	}

	data.ID = input.Module.ID
	data.Name = input.Module.Name
	data.Description = input.Module.Description
	data.State = input.Module.State

	questions := map[string]interface{}{}
	if input.Module.Questions != nil {
		json.Unmarshal(input.Module.Questions, &questions)
	}
	data.Questions = questions

	data.ModuleCategory = Resource{
		ID:   input.ModuleCategory.ID,
		Name: input.ModuleCategory.Name,
	}

	data.Creator = User{
		ID: input.Module.CreatorID,
	}
	data.Creator.Fill(ctx)

	data.CreatedAt = input.Module.CreatedAt
	data.UpdatedAt = input.Module.UpdatedAt

	return *data
}

func (data *Module) ExtractByViewProto(ctx context.Context, input *pb.ModuleJoinModuleAclInfo) Module {
	if input == nil {
		return Module{}
	}

	data.ID = input.ModuleID
	data.Name = input.ModuleName
	data.Description = input.ModuleDescription
	data.State = input.State

	questions := map[string]interface{}{}
	if input.Questions != nil {
		json.Unmarshal(input.Questions, &questions)
	}
	data.Questions = questions

	data.ModuleCategory = Resource{
		ID:   input.ModuleCategoryID,
		Name: input.ModuleCategoryName,
	}

	data.Creator = User{
		ID: input.ModuleCreatorID,
	}
	data.Creator.Fill(ctx)

	data.CreatedAt = input.ModuleCreatedAt
	data.UpdatedAt = input.ModuleUpdatedAt

	return *data
}
