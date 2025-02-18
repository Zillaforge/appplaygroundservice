package keystone

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/openstack/common"
	"AppPlaygroundService/utility"
	"context"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/applicationcredentials"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (n *Keystone) CreateAppCredential(ctx context.Context, input *common.CreateAppCredentialInput) (output *common.CreateAppCredentialOutput, err error) {
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

	createOpts := applicationcredentials.CreateOpts{
		Name:         input.Name,
		Description:  input.Description,
		Unrestricted: input.Unrestricted,
	}

	appCred, err := applicationcredentials.Create(n.sc, input.OpstkUserID, createOpts).Extract()
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "applicationcredentials.Create(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("createOpts", createOpts),
		).Error(err.Error())
		return nil, err
	}

	output = &common.CreateAppCredentialOutput{
		ID:             appCred.ID,
		Name:           appCred.Name,
		Secret:         appCred.Secret,
		OpstkProjectID: appCred.ProjectID,
	}
	return
}
