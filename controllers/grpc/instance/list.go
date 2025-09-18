package instance

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

func (m *Method) ListInstances(ctx context.Context, input *pb.ListInput) (output *pb.ListInstancesOutput, err error) {
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

	whereInput := storCom.ListInstancesWhere{}
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
		return &pb.ListInstancesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listInstancesInput := &storCom.ListInstancesInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listInstancesOutput, err := storages.Use().ListInstances(ctx, listInstancesInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListInstances(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInstancesInput),
		).Error(err.Error())
		return &pb.ListInstancesOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListInstancesOutput{
		Count: listInstancesOutput.Count,
	}
	for _, Instance := range listInstancesOutput.Instances {
		output.Data = append(output.Data, m.storage2pb(&Instance))
	}
	return
}
