package execution

import (
	"AppPlaygroundService/modules/openstack/keystone"
	"AppPlaygroundService/modules/openstack/neutron"

	"github.com/gophercloud/gophercloud"
)

type AdminResource interface {
	Neutron() *neutron.Neutron
	Keystone() *keystone.Keystone
}

type Admin struct {
	providerClient func(*string, *string, *string) (*gophercloud.ProviderClient, error)
	serviceClient  func(*gophercloud.ProviderClient, string) (*gophercloud.ServiceClient, error)

	neutron  func() *neutron.Neutron
	keystone func() *keystone.Keystone
}

// Neutron ...
func (p *Admin) Neutron() *neutron.Neutron {
	return p.neutron()
}

// Keystone ...
func (p *Admin) Keystone() *keystone.Keystone {
	return p.keystone()
}

func (cfg *Connection) Admin() AdminResource {
	return &Admin{
		providerClient: cfg.providerClient,
		serviceClient:  cfg.serviceClient,
		neutron: func() *neutron.Neutron {
			return cfg.Neutron(cfg.AdminProject, "")
		},
		keystone: func() *keystone.Keystone {
			return cfg.Keystone(cfg.AdminProject, "")
		},
	}
}
