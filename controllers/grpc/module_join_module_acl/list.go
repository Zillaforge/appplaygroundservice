package modulejoinmoduleacl

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

func (m *Method) ListModuleJoinModuleAcls(ctx context.Context, input *pb.ListModuleJoinModuleAclsInput) (output *pb.ListModuleJoinModuleAclsOutput, err error) {
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

	whereInput := storCom.ListModuleJoinModuleAclsWhere{}
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
		return &pb.ListModuleJoinModuleAclsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listModuleJoinModuleAclsInput := &storCom.ListModuleJoinModuleAclsInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
		ProjectID:  input.ProjectID,
	}
	listModuleJoinModuleAclsOutput, err := storages.Use().ListModuleJoinModuleAcls(ctx, listModuleJoinModuleAclsInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListModuleJoinModuleAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModuleJoinModuleAclsInput),
		).Error(err.Error())
		return &pb.ListModuleJoinModuleAclsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListModuleJoinModuleAclsOutput{
		Count: listModuleJoinModuleAclsOutput.Count,
	}
	for _, moduleJoinModuleAcl := range listModuleJoinModuleAclsOutput.ModuleJoinModuleAcls {
		data := m.storage2pb(&moduleJoinModuleAcl)

		if input.Language != "" {
			var getQuestionsErr error
			data.Questions, getQuestionsErr = m.getQuestions(moduleJoinModuleAcl.ModuleID, input.Language)
			if getQuestionsErr != nil {
				getQuestionsErr = tkErr.New(cCnt.GRPCGetModuleQuestionsFailedErr).WithInner(getQuestionsErr)
				zap.L().With(
					zap.String(cnt.Controller, "getQuestions(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.String("moduleID", moduleJoinModuleAcl.ModuleID),
					zap.String("language", input.Language),
				).Error(getQuestionsErr.Error())
			}
		}

		output.Data = append(output.Data, data)
	}
	return
}
