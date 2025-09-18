package keystone

import (
	cnt "AppPlaygroundService/constants"

	"github.com/gophercloud/gophercloud"
	tkErr "github.com/Zillaforge/toolkits/errors"
)

type Keystone struct {
	namespace string
	projectID string
	sc        *gophercloud.ServiceClient
}

func New(namespace, projectID string, sc *gophercloud.ServiceClient) *Keystone {
	return &Keystone{
		namespace: namespace,
		projectID: projectID,
		sc:        sc,
	}
}

func (n *Keystone) SetServiceClient(namespace string, sc *gophercloud.ServiceClient) *Keystone {
	n.namespace = namespace
	n.sc = sc
	return n
}

func (n *Keystone) checkConnection() error {
	if n.sc == nil {
		return tkErr.New(cnt.OpenstackConnectionIsNotCreatedErr)
	}
	return nil
}
