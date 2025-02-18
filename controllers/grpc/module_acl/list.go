package moduleacl

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/grpc"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (m *Method) ListModuleAcls(ctx context.Context, input *pb.ListInput) (output *pb.ListModuleAclsOutput, err error) {
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

	whereInput := storCom.ListModuleAclsWhere{}
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
		return &pb.ListModuleAclsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listModuleAclsInput := &storCom.ListModuleAclsInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listModuleAclsOutput, err := storages.Use().ListModuleAcls(ctx, listModuleAclsInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListModuleAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModuleAclsInput),
		).Error(err.Error())
		return &pb.ListModuleAclsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListModuleAclsOutput{
		Count: listModuleAclsOutput.Count,
		Data:  []*pb.ModuleAclInfo{},
	}
	for _, moduleAcl := range listModuleAclsOutput.ModuleAcls {
		output.Data = append(output.Data, m.storage2pb(&moduleAcl))
	}
	return
}
