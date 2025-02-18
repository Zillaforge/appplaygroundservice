package app_credential

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/openstack"
	opstkCom "AppPlaygroundService/modules/openstack/common"
	"AppPlaygroundService/modules/opstkidentity"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (m *Method) CreateAppCredential(ctx context.Context, input *pb.AppCredentialInfo) (output *pb.AppCredentialInfo, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	opstkUID, err := opstkidentity.Use().GetOpstkUID(ctx, input.UserID)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "o.GetOpstkUID"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input.UserID),
		).Error(err.Error())
		return
	}

	// Create app credential in openstack
	opstkInput := &opstkCom.CreateAppCredentialInput{
		Name:         input.Name,
		Description:  "Created by APS",
		Unrestricted: false,
		OpstkUserID:  opstkUID,
	}
	opstkOutput, opstkErr := openstack.Namespace(input.Namespace).Keystone(input.ProjectID, input.UserID).CreateAppCredential(ctx, opstkInput)
	if opstkErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "openstack.Namespace(...).Keystone(...).CreateAppCredential"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", opstkInput),
		).Error(opstkErr.Error())
		return
	}

	// Create app credential in storage
	createInput := &storCom.CreateAppCredentialInput{
		AppCredential: tables.AppCredential{
			UserID:    input.UserID,
			ProjectID: input.ProjectID,
			ID:        opstkOutput.ID,
			Name:      opstkOutput.Name,
			Secret:    opstkOutput.Secret,
			Namespace: input.Namespace,
		},
	}
	createOutput, err := storages.Use().CreateAppCredential(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageAppCredentialExistErrCode:
				err = tkErr.New(cCnt.GRPCAppCredentialExistErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateAppCredential()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&createOutput.AppCredential)
	return
}
