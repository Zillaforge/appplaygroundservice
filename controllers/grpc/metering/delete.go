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

func (m *Method) DeleteMetering(ctx context.Context, input *pb.DeleteMeteringInput) (output *pb.DeleteMeteringOutput, err error) {
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

	deleteMeteringInput := &storCom.DeleteMeteringInput{
		ID: input.ApplicationID,
	}
	if _, deleteMeteringErr := storages.Use().DeleteMetering(ctx, deleteMeteringInput); deleteMeteringErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().DeleteMetering()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteMeteringInput),
		).Error(deleteMeteringErr.Error())
		return nil, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	return nil, nil
}
