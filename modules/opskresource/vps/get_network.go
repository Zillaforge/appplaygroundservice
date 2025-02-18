package vps

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/opskresource/common"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualplatformserviceclient/pb"
)

func (h *Handler) GetNetwork(ctx context.Context, input *common.GetNetworkInput) (output *common.GetNetworkOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"error":  &err,
	})

	getNetworkInput := &pb.IDInput{
		ID: input.ID,
	}
	getNetworkOutput, err := h.poolHandler.Network().Get(getNetworkInput, ctx)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case StorageRecordNotFoundErrCode:
				err = tkErr.New(cnt.OpskResourceRecordNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.OpskResource, "h.poolHandler.Network().Get(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getNetworkInput),
		).Error(err.Error())
		return
	}

	output = &common.GetNetworkOutput{
		ID:        getNetworkOutput.ID,
		Name:      getNetworkOutput.Name,
		ProjectID: getNetworkOutput.ProjectID,
		Namespace: getNetworkOutput.Namespace,
		RouterID:  getNetworkOutput.RouterID,
		SubnetID:  getNetworkOutput.SubnetID,
	}
	return output, nil
}
