package execution

import (
	cnt "AppPlaygroundService/constants"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

type Connection struct {
	IdentityEndpoint  string
	AdminUsername     string
	AdminPassword     string
	AdminProject      string
	DomainName        string
	AllowReauth       bool
	PidSource         string
	Pid               func(string) string
	UsernameSource    string
	PasswordOTPSecret string
	Username          func(string) string
	Password          func(string) string

	namespace string
}

func New(p *Connection) (*Connection, error) {
	_, err := p.providerClient(nil, &p.AdminUsername, &p.AdminPassword)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (cfg *Connection) providerClient(projectID *string, username *string, password *string) (*gophercloud.ProviderClient, error) {
	authOptions := gophercloud.AuthOptions{
		IdentityEndpoint: cfg.IdentityEndpoint,
		Username:         *username,
		Password:         *password,
		DomainName:       cfg.DomainName,
		AllowReauth:      cfg.AllowReauth,
	}
	if projectID != nil {
		authOptions.Scope = &gophercloud.AuthScope{
			ProjectID: *projectID,
		}
	}
	return openstack.AuthenticatedClient(authOptions)
}

func (Connection) serviceClient(pc *gophercloud.ProviderClient, key string) (*gophercloud.ServiceClient, error) {
	var scFunc func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	switch key {
	case _neutronResource:
		scFunc = openstack.NewNetworkV2
	case _keystoneResource:
		scFunc = openstack.NewIdentityV3
	default:
		return nil, tkErr.New(cnt.OpenstackTypeIsNotSupportedErr)
	}
	return scFunc(pc, gophercloud.EndpointOpts{})
}

func (cfg *Connection) SetNamespace(namespace string) {
	cfg.namespace = namespace
}
