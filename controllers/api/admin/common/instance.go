package common

import (
	"encoding/json"

	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

type Instance struct {
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

type InstanceInterface interface {
	ExtractByProto(proto *pb.InstanceDetail)
}

func (i *Instance) ExtractByProto(proto *pb.InstanceDetail) {
	if proto == nil {
		return
	}
	i.ID = proto.Instance.ID
	i.Name = proto.Instance.Name
	i.Application = Resource{
		ID:   proto.Application.ID,
		Name: proto.Application.Name,
	}
	i.ProjectID = proto.Instance.ProjectID
	i.ReferenceID = proto.Instance.ReferenceID

	metadata := map[string]interface{}{}
	json.Unmarshal(proto.Instance.Extra, &metadata)
	i.Extra = metadata

	i.FloatingIPID = proto.Instance.FloatingIPID
	i.FloatingIPAddress = proto.Instance.FloatingIPAddress

	i.CreatedAt = proto.Instance.CreatedAt
	i.UpdatedAt = proto.Instance.UpdatedAt
}
