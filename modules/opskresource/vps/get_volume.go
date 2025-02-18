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

func (h *Handler) GetVolume(ctx context.Context, input *common.GetVolumeInput) (output *common.GetVolumeOutput, err error) {
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

	getVolumeInput := &pb.IDInput{
		ID: input.ID,
	}
	getVolumeOutput, err := h.poolHandler.Volume().Get(getVolumeInput, ctx)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case StorageRecordNotFoundErrCode:
				err = tkErr.New(cnt.OpskResourceRecordNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.OpskResource, "h.poolHandler.Volume().Get(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getVolumeInput),
		).Error(err.Error())
		return
	}

	output = &common.GetVolumeOutput{
		ID:        getVolumeOutput.ID,
		Name:      getVolumeOutput.Name,
		ProjectID: getVolumeOutput.Project.ID,
		UserID:    getVolumeOutput.User.ID,
		Namespace: getVolumeOutput.Namespace,
		Status:    getVolumeOutput.Status,
		Type:      getVolumeOutput.Type,
		Size:      getVolumeOutput.Size,
	}
	return output, nil
}
