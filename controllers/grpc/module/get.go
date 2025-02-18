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

func (m *Method) GetModule(ctx context.Context, input *pb.GetWithLanguageInput) (output *pb.ModuleDetail, err error) {
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

	getInput := &storCom.GetModuleInput{
		ID: input.ID,
	}

	getOutput, err := storages.Use().GetModule(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageModuleNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCModuleNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetModule()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&getOutput.Module)

	if input.Language != "" {
		var getQuestionsErr error
		output.Module.Questions, getQuestionsErr = m.getQuestions(getOutput.Module.ID, input.Language)
		if getQuestionsErr != nil {
			err = tkErr.New(cCnt.GRPCGetModuleQuestionsFailedErr).WithInner(getQuestionsErr)
			zap.L().With(
				zap.String(cnt.GRPC, "getQuestions(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.String("moduleID", input.ID),
				zap.String("language", input.Language),
			).Error(err.Error())
			return
		}
	}

	return
}
