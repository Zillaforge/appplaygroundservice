package common

import (
	"encoding/json"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

type (
	Instance struct {
		ID                string                 `json:"id"`
		Name              string                 `json:"name"`
		Application       Resource               `json:"application"`
		ProjectID         string                 `json:"projectId"`
		ReferenceID       string                 `json:"referenceId"`
		Extra             map[string]interface{} `json:"extra"`
		FloatingIPID      string                 `json:"floatingIpId"`
		FloatingIPAddress string                 `json:"floatingIpAddress"`
		CreatedAt         string                 `json:"createdAt"`
		UpdatedAt         string                 `json:"updatedAt"`
		_                 struct{}
	}
)

func (data *Instance) ExtractByProto(input *pb.InstanceDetail) Instance {
	data.ID = input.Instance.ID
	data.Name = input.Instance.Name

	data.Application = Resource{
		ID:   input.Application.ID,
		Name: input.Application.Name,
	}

	data.ProjectID = input.Instance.ProjectID
	data.ReferenceID = input.Instance.ReferenceID

	extra := map[string]interface{}{}
	if input.Instance.Extra != nil {
		json.Unmarshal(input.Instance.Extra, &extra)
	}
	data.Extra = extra

	data.FloatingIPID = input.Instance.FloatingIPID
	data.FloatingIPAddress = input.Instance.FloatingIPAddress

	data.CreatedAt = input.Instance.CreatedAt
	data.UpdatedAt = input.Instance.UpdatedAt

	return *data
}
