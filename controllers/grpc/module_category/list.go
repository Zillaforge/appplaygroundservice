package modulecategory

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

func (m *Method) ListModuleCategories(ctx context.Context, input *pb.ListInput) (output *pb.ListModuleCategoriesOutput, err error) {
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

	whereInput := storCom.ListModuleCategoriesWhere{}
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
		return &pb.ListModuleCategoriesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listModuleCategoriesInput := &storCom.ListModuleCategoriesInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listModuleCategoriesOutput, err := storages.Use().ListModuleCategories(ctx, listModuleCategoriesInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListModuleCategories(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModuleCategoriesInput),
		).Error(err.Error())
		return &pb.ListModuleCategoriesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListModuleCategoriesOutput{
		Count: listModuleCategoriesOutput.Count,
	}
	for _, ModuleCategory := range listModuleCategoriesOutput.ModuleCategories {
		output.Data = append(output.Data, m.storage2pb(&ModuleCategory))
	}
	return
}
