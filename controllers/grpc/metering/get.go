package metering

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

func (m *Method) GetMetering(ctx context.Context, input *pb.GetInput) (output *pb.MeteringInfo, err error) {
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

	getInput := &storCom.GetMeteringInput{
		ID: input.ID,
	}

	getOutput, err := storages.Use().GetMetering(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageMeteringNotFoundErr.Code():
				return nil, tkErr.New(cCnt.GRPCMeteringNotFoundErr)
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetMetering()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		return nil, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	return m.storage2pb(&getOutput.Metering), nil
}
