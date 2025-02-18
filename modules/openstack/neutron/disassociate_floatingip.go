package neutron

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/openstack/common"
	"AppPlaygroundService/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (n *Neutron) DisassociateFloatingIp(ctx context.Context, input *common.DisassociateFloatingIpInput) (err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input": &input,
			"error": &err,
		},
	)

	if err = n.checkConnection(); err != nil {
		zap.L().With(
			zap.String(cnt.Module, "n.checkConnection()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", n.namespace),
		).Error(err.Error())
		return
	}

	updateOpts := floatingips.UpdateOpts{
		PortID: nil,
	}

	_, err = floatingips.Update(n.sc, input.FloatingIpID, updateOpts).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "floatingips.Update.Extract()"),
			zap.String(cnt.RequestID, requestID),
			zap.String("namespace", n.namespace),
			zap.String("fip-id", input.FloatingIpID),
		).Error(err.Error())
		err = tkErr.New(cnt.OpenstackInternalServerErr).WithInner(err)
		return
	}
	return
}
