package app_credential

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages"
	"AppPlaygroundService/utility"
	"context"
	"fmt"

	storCom "AppPlaygroundService/storages/common"

	"go.uber.org/zap"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) GetAppCredential(ctx context.Context, input *pb.GetAppCredentialInput) (output *pb.AppCredentialInfo, err error) {
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

	getInput := &storCom.GetAppCredentialInput{
		UserID:    &input.UserID,
		ProjectID: &input.ProjectID,
		Namespace: &input.Namespace,
	}

	getOutput, err := storages.Use().GetAppCredential(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageAppCredentialNotFoundErrCode:
				// If app credential not found, create a new
				createInput := &pb.AppCredentialInfo{
					UserID:    input.UserID,
					ProjectID: input.ProjectID,
					Namespace: input.Namespace,
					Name:      fmt.Sprintf("aps-%s-%s-%s", input.UserID, input.ProjectID, tkUtils.GetRandAlphanumericLower(8)),
				}
				createOutput, createErr := m.CreateAppCredential(ctx, createInput)
				if createErr != nil {
					zap.L().With(
						zap.String(cnt.GRPC, "m.CreateAppCredential(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.Any("input", createInput),
					).Error(createErr.Error())
					err = createErr
					return
				}
				output = createOutput
				return output, nil
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetAppCredential()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&getOutput.AppCredential)
	return
}
