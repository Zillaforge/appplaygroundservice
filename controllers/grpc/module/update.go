package module

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

func (m *Method) UpdateModule(ctx context.Context, input *pb.UpdateModuleInput) (output *pb.ModuleDetail, err error) {
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

	updateInput := &storCom.UpdateModuleInput{
		ID: input.ID,
		UpdateData: &storCom.ModuleUpdateInfo{
			Name:        input.Name,
			Description: input.Description,
			State:       input.State,
			Public:      input.Public,
		},
	}

	updateOutput, err := storages.Use().UpdateModule(ctx, updateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageModuleExistErrCode:
				err = tkErr.New(cCnt.GRPCModuleExistErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateModule()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}
	output = m.storage2pb(&updateOutput.Module)

	if input.Language != "" {
		var getQuestionsErr error
		output.Module.Questions, getQuestionsErr = m.getQuestions(input.ID, input.Language)
		if getQuestionsErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "getQuestions(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("moduleID", input.ID),
				zap.String("language", input.Language),
			).Warn(getQuestionsErr.Error())
		}
	}
	return
}
