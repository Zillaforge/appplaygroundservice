package common

import (
	"context"
	"encoding/json"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type (
	Application struct {
		ID             string                 `json:"id"`
		Name           string                 `json:"name"`
		Description    string                 `json:"description"`
		State          string                 `json:"state"`
		Answers        map[string]interface{} `json:"answers"`
		Module         Resource               `json:"module"`
		ModuleCategory Resource               `json:"moduleCategory"`
		Namespace      string                 `json:"namespace"`
		Shiftable      bool                   `json:"shiftable"`
		ProjectID      string                 `json:"projectId"`
		Creator        User                   `json:"creator"`
		Updater        User                   `json:"updater"`
		Extra          map[string]interface{} `json:"extra"`
		CreatedAt      string                 `json:"createdAt"`
		UpdatedAt      string                 `json:"updatedAt"`
		_              struct{}
	}
)

func (data *Application) ExtractByProto(ctx context.Context, input *pb.ApplicationDetail) Application {
	var (
		funcName = tkUtils.NameOfFunction().String()
		err      error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	data.ID = input.Application.ID
	data.Name = input.Application.Name
	data.Description = input.Application.Description
	data.State = input.Application.State

	answers := map[string]interface{}{}
	if input.Application.Answers != nil {
		json.Unmarshal(input.Application.Answers, &answers)
	}
	data.Answers = answers

	data.Module = Resource{
		ID:   input.Module.ID,
		Name: input.Module.Name,
	}

	data.ModuleCategory = Resource{
		ID:   input.ModuleCategory.ID,
		Name: input.ModuleCategory.Name,
	}

	data.Namespace = input.Application.Namespace
	data.Shiftable = input.Application.Shiftable
	data.ProjectID = input.Application.ProjectID

	data.Creator = User{
		ID: input.Application.CreatorID,
	}
	data.Creator.Fill(ctx)

	data.Updater = User{
		ID: input.Application.UpdaterID,
	}
	data.Updater.Fill(ctx)

	extra := map[string]interface{}{}
	if input.Application.Extra != nil {
		json.Unmarshal(input.Application.Extra, &extra)
	}
	data.Extra = extra

	data.CreatedAt = input.Application.CreatedAt
	data.UpdatedAt = input.Application.UpdatedAt

	return *data
}
