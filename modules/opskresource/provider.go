package opskresource

import (
	"AppPlaygroundService/modules/opskresource/common"
	"AppPlaygroundService/modules/opskresource/vps"
	"AppPlaygroundService/services"
	"context"
	"fmt"

	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/mviper"
)

const (
	_vpsType = "vps"
)

var _provider Provider

// Provider define functions which can be implemented for openstack resource
type Provider interface {
	// FloatingIP
	GetFloatingIP(ctx context.Context, input *common.GetFloatingIPInput) (output *common.GetFloatingIPOutput, err error)
	UpdateFloatingIPStatus(ctx context.Context, input *common.UpdateFloatingIPStatusInput) (output *common.UpdateFloatingIPStatusOutput, err error)

	// SecurityGroup
	GetSecurityGroup(ctx context.Context, input *common.GetSecurityGroupInput) (output *common.GetSecurityGroupOutput, err error)

	// Flavor
	GetFlavor(ctx context.Context, input *common.GetFlavorInput) (output *common.GetFlavorOutput, err error)

	// Network
	GetNetwork(ctx context.Context, input *common.GetNetworkInput) (output *common.GetNetworkOutput, err error)

	// Keypair
	GetKeypair(ctx context.Context, input *common.GetKeypairInput) (output *common.GetKeypairOutput, err error)

	// Volume
	GetVolume(ctx context.Context, input *common.GetVolumeInput) (output *common.GetVolumeOutput, err error)
}

// New ...
func New(service string) {
	kind := services.ServiceMap[service].Kind
	switch kind {
	case _vpsType:
		zap.L().Info(fmt.Sprintf("openstack resource is %s(%s) service mode", kind, service))
		_provider = vps.New(vps.WithServiceName(mviper.GetString("openstack_resource.service")))
	default:
		panic(fmt.Errorf("openstack resource does not support %s mode", kind))
	}
}

// Use returns openstack resource instance
func Use() Provider {
	return _provider
}

// Replace replaces global provider by p
func Replace(p Provider) {
	_provider = p
}
