package common

import (
	auth "AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	"context"
	"encoding/json"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

type Module struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	State          string                 `json:"state"`
	Questions      map[string]interface{} `json:"questions"`
	Location       string                 `json:"location"`
	Public         bool                   `json:"public"`
	ModuleCategory Resource               `json:"moduleCategory"`
	Creator        User                   `json:"creator"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
}

type ModuleInterface interface {
	ExtractByProto(ctx context.Context, proto *pb.ModuleDetail)
}

func (m *Module) ExtractByProto(ctx context.Context, proto *pb.ModuleDetail) {
	if proto == nil {
		return
	}

	m.ID = proto.Module.ID
	m.Name = proto.Module.Name
	m.Description = proto.Module.Description
	m.State = proto.Module.State

	metadata := map[string]interface{}{}
	json.Unmarshal(proto.Module.Questions, &metadata)
	m.Questions = metadata

	m.Location = proto.Module.Location
	m.Public = proto.Module.Public

	m.ModuleCategory.ID = proto.ModuleCategory.ID
	m.ModuleCategory.Name = proto.ModuleCategory.Name

	// 用 CreatorID 換 User 資訊
	getUserInput := &authCom.GetUserInput{
		ID:        proto.Module.CreatorID,
		Cacheable: true,
	}
	getCreatorOutput, getCreatorErr := auth.Use().GetUser(ctx, getUserInput)
	if getCreatorErr == nil {
		m.Creator = User{
			ID:      getCreatorOutput.ID,
			Account: getCreatorOutput.Account,
			Name:    getCreatorOutput.DisplayName,
		}
	}

	m.CreatedAt = proto.Module.CreatedAt
	m.UpdatedAt = proto.Module.UpdatedAt
}
