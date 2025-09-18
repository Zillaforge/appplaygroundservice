package neutron

import (
	cnt "AppPlaygroundService/constants"

	"github.com/gophercloud/gophercloud"
	tkErr "github.com/Zillaforge/toolkits/errors"
)

type Neutron struct {
	namespace string
	projectID string
	sc        *gophercloud.ServiceClient
}

func New(namespace, projectID string, sc *gophercloud.ServiceClient) *Neutron {
	return &Neutron{
		namespace: namespace,
		projectID: projectID,
		sc:        sc,
	}
}

func (n *Neutron) SetServiceClient(namespace string, sc *gophercloud.ServiceClient) *Neutron {
	n.namespace = namespace
	n.sc = sc
	return n
}

func (n *Neutron) checkConnection() error {
	if n.sc == nil {
		return tkErr.New(cnt.OpenstackConnectionIsNotCreatedErr)
	}
	return nil
}
