package common

import (
	auth "AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	"context"
	"encoding/json"

	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

type Application struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	State          string                 `json:"state"`
	Answers        map[string]interface{} `json:"answers"`
	Namespace      string                 `json:"namespace"`
	Shiftable      bool                   `json:"shiftable"`
	ProjectID      string                 `json:"projectId"`
	Creator        User                   `json:"creator"`
	Updater        *User                  `json:"updater"`
	Module         Resource               `json:"module"`
	ModuleCategory Resource               `json:"moduleCategory"`
	Extra          map[string]interface{} `json:"extra"`
	CreatedAt      string                 `json:"createdAt"`
	UpdatedAt      string                 `json:"updatedAt"`
	_              struct{}
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
}

type Resource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ApplicationInterface interface {
	ExtractByProto(ctx context.Context, proto *pb.ApplicationDetail)
}

func (a *Application) ExtractByProto(ctx context.Context, proto *pb.ApplicationDetail) {
	if proto == nil {
		return
	}
	if proto.Application != nil {
		a.ID = proto.Application.ID
		a.Name = proto.Application.Name
		a.Description = proto.Application.Description
		a.State = proto.Application.State

		metadata := map[string]interface{}{}
		json.Unmarshal(proto.Application.Answers, &metadata)
		a.Answers = metadata

		a.Namespace = proto.Application.Namespace
		a.Shiftable = proto.Application.Shiftable
		a.ProjectID = proto.Application.ProjectID

		// 用 CreatorID 換 User 資訊
		getUserInput := &authCom.GetUserInput{
			ID:        proto.Application.CreatorID,
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

		// 用 UpdaterID 換 User 資訊
		if proto.Application.UpdaterID != "" {
			getUpdaterInput := &authCom.GetUserInput{
				ID:        proto.Application.UpdaterID,
				Cacheable: true,
			}
			getUpdaterOutput, getUpdaterErr := auth.Use().GetUser(ctx, getUpdaterInput)
			if getUpdaterErr == nil {
				a.Creator = User{
					ID:      getUpdaterOutput.ID,
					Account: getUpdaterOutput.Account,
					Name:    getUpdaterOutput.DisplayName,
				}
			}
		}

		a.Module = Resource{
			ID:   proto.Module.ID,
			Name: proto.Module.Name,
		}

		a.ModuleCategory = Resource{
			ID:   proto.ModuleCategory.ID,
			Name: proto.ModuleCategory.Name,
		}

		extra := map[string]interface{}{}
		json.Unmarshal(proto.Application.Extra, &extra)
		a.Extra = extra

		a.CreatedAt = proto.Application.CreatedAt
		a.UpdatedAt = proto.Application.UpdatedAt
	}
}
