package keystone

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/openstack/common"
	"AppPlaygroundService/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/applicationcredentials"
	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (n *Keystone) DeleteAppCredential(ctx context.Context, input *common.DeleteAppCredentialInput) (err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	_, f := tracer.StartWithContext(ctx, funcName)
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

	result := applicationcredentials.Delete(n.sc, input.OpstkUserID, input.ID)
	if result.Err != nil {
		zap.L().With(
			zap.String(cnt.Module, "applicationcredentials.Delete(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String("userID", input.OpstkUserID),
			zap.String("id", input.ID),
		).Error(result.Err.Error())
		return result.Err
	}
	return
}
