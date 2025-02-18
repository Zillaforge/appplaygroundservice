package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/application"
	"AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (m *Method) GetAppLogs(ctx context.Context, input *pb.GetAppLogsInput) (output *pb.GetAppLogsOutput, err error) {
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

	output = &pb.GetAppLogsOutput{
		Logs: "",
	}

	// get Application's projectID
	getInput := &storCom.GetApplicationInput{
		ID: input.ApplicationID,
	}
	getOutput, err := storages.Use().GetApplication(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageApplicationNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCApplicationNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetApplication()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	getLogsInput := common.GetLogsInput{
		ApplicationID: input.ApplicationID,
		ProjectID:     getOutput.Application.ProjectID,
	}
	getLogsOutput, err := application.Use().GetLogs(ctx, getLogsInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "application.Use().GetLogs(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getLogsInput),
		).Error(err.Error())
		return
	}

	output.Logs = getLogsOutput.Logs

	return
}
