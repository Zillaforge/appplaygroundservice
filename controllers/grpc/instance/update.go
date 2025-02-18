package instance

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

func (m *Method) UpdateInstance(ctx context.Context, input *pb.UpdateInstanceInput) (output *pb.InstanceDetail, err error) {
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

	updateInput := &storCom.UpdateInstanceInput{
		ID: input.ID,
		UpdateData: &storCom.InstanceUpdateInfo{
			Name:  input.Name,
			Extra: input.Extra,
		},
	}
	_, err = storages.Use().UpdateInstance(ctx, updateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageInstanceNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCInstanceNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateInstance()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	// get updated instance
	getInstanceInput := &storCom.GetInstanceInput{
		ID: input.ID,
	}
	getInstanceOutput, err := storages.Use().GetInstance(ctx, getInstanceInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageInstanceNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCInstanceNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetInstance()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInstanceInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&getInstanceOutput.Instance)
	return
}
