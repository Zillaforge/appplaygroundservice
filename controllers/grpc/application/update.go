package application

import (
	cnt "AppPlaygroundService/constants"
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

func (m *Method) UpdateApplication(ctx context.Context, input *pb.UpdateApplicationInput) (output *pb.ApplicationDetail, err error) {
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

	updateInput := &storCom.UpdateApplicationInput{
		ID: input.ID,
		UpdateData: &storCom.ApplicationUpdateInfo{
			Name:        input.Name,
			Description: input.Description,
			State:       input.State,
			Answers:     input.Answers,
			Namespace:   input.Namespace,
			Shiftable:   input.Shiftable,
			UpdaterID:   input.UpdaterID,
			Extra:       input.Extra,
		},
	}

	updateOutput, err := storages.Use().UpdateApplication(ctx, updateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageApplicationNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCApplicationNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateApplication()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	if input.Language != "" {
		updateOutput.Application.Answers, err = m.getAnswers(
			string(updateOutput.Application.Answers),
			updateOutput.Application.ModuleID,
			input.Language)
		if err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "m.getAnswers()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("ansStr", string(updateOutput.Application.Answers)),
				zap.String("language", input.Language),
			).Error(err.Error())
			return
		}
	}

	output = m.storage2pb(&updateOutput.Application)
	return
}
