package execution

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/openstack/keystone"
	"AppPlaygroundService/modules/openstack/neutron"

	"go.uber.org/zap"
)

const (
	_neutronResource  = "neutron"
	_keystoneResource = "keystone"
)

// Neutron ...
// If userID is not provided, it will use the default admin user
func (p *Connection) Neutron(projectID string, userID string) *neutron.Neutron {
	pid := p.Pid(projectID)
	username := p.Username(userID)
	password := p.Password(userID)
	pc, err := p.providerClient(&pid, &username, &password)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.providerClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &neutron.Neutron{}
	}
	sc, err := p.serviceClient(pc, _neutronResource)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.serviceClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &neutron.Neutron{}
	}

	return neutron.New(p.namespace, pid, sc)
}

// Neutron ...
// If userID is not provided, it will use the default admin user
func (p *Connection) Keystone(projectID string, userID string) *keystone.Keystone {
	pid := p.Pid(projectID)
	username := p.Username(userID)
	password := p.Password(userID)
	pc, err := p.providerClient(&pid, &username, &password)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.providerClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &keystone.Keystone{}
	}
	sc, err := p.serviceClient(pc, _keystoneResource)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.serviceClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &keystone.Keystone{}
	}

	return keystone.New(p.namespace, pid, sc)
}
