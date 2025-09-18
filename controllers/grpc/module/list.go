package module

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/grpc"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) ListModules(ctx context.Context, input *pb.ListWithLanguageInput) (output *pb.ListModulesOutput, err error) {
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

	whereInput := storCom.ListModulesWhere{}
	if err := querydecoder.ShouldBindWhereSlice(&whereInput, input.Where); err != nil {
		if e, ok := tkErr.IsError(grpc.WhereErrorParser(err)); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				return output, e
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "querydecoder.ShouldBindWhereSlice(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input.where", input.Where),
		).Error(err.Error())
		return &pb.ListModulesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listModulesInput := &storCom.ListModulesInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listModulesOutput, err := storages.Use().ListModules(ctx, listModulesInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListModules(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModulesInput),
		).Error(err.Error())
		return &pb.ListModulesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListModulesOutput{
		Count: listModulesOutput.Count,
	}
	for _, module := range listModulesOutput.Modules {
		data := m.storage2pb(&module)

		if input.Language != "" {
			var getQuestionsErr error
			data.Module.Questions, getQuestionsErr = m.getQuestions(module.ID, input.Language)
			if getQuestionsErr != nil {
				getQuestionsErr = tkErr.New(cCnt.GRPCGetModuleQuestionsFailedErr).WithInner(getQuestionsErr)
				zap.L().With(
					zap.String(cnt.GRPC, "getQuestions(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.String("moduleID", module.ID),
					zap.String("language", input.Language),
				).Error(getQuestionsErr.Error())
			}
		}

		output.Data = append(output.Data, data)
	}
	return
}
