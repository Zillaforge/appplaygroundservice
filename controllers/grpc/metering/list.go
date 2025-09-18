package metering

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) ListMeterings(ctx context.Context, input *pb.ListInput) (output *pb.ListMeteringsOutput, err error) {
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

	listMeteringsInput := &storCom.ListMeteringsInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
	}
	listMeteringsOutput, listMeteringsErr := storages.Use().ListMeterings(ctx, listMeteringsInput)
	if listMeteringsErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListMeterings(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listMeteringsInput),
		).Error(listMeteringsErr.Error())
		return &pb.ListMeteringsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	countMeteringsOutput, countMeteringsErr := storages.Use().CountMetering(ctx, &storCom.CountMeteringInput{})
	if countMeteringsErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CountMetering(...)"),
			zap.String(cnt.RequestID, requestID),
		).Error(countMeteringsErr.Error())
		return &pb.ListMeteringsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListMeteringsOutput{
		Count: countMeteringsOutput.Count,
	}
	for _, metering := range listMeteringsOutput.Meterings {
		output.Data = append(output.Data, m.storage2pb(&metering))
	}
	return output, nil
}
