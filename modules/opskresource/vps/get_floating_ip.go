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

func (h *Handler) GetFloatingIP(ctx context.Context, input *common.GetFloatingIPInput) (output *common.GetFloatingIPOutput, err error) {
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

	getFloatingIPInput := &pb.IDInput{
		ID: input.ID,
	}
	getFloatingIPOutput, err := h.poolHandler.FloatingIP().Get(getFloatingIPInput, ctx)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case StorageRecordNotFoundErrCode:
				err = tkErr.New(cnt.OpskResourceRecordNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.OpskResource, "h.poolHandler.Flavor().Get(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getFloatingIPInput),
		).Error(err.Error())
		return
	}

	output = &common.GetFloatingIPOutput{
		ID:         getFloatingIPOutput.ID,
		UUID:       getFloatingIPOutput.UUID,
		Name:       getFloatingIPOutput.Name,
		ProjectID:  getFloatingIPOutput.ProjectID,
		UserID:     getFloatingIPOutput.UserID,
		Namespace:  getFloatingIPOutput.Namespace,
		Status:     getFloatingIPOutput.Status,
		Reserved:   getFloatingIPOutput.Reserved,
		DeviceType: getFloatingIPOutput.DeviceType,
		DeviceID:   getFloatingIPOutput.DeviceID,
		Address:    getFloatingIPOutput.Address,
	}
	return output, nil
}
