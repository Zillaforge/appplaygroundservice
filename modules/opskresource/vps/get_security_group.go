package vps

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/opskresource/common"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	"github.com/Zillaforge/virtualplatformserviceclient/pb"
)

func (h *Handler) GetSecurityGroup(ctx context.Context, input *common.GetSecurityGroupInput) (output *common.GetSecurityGroupOutput, err error) {
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

	getSgInput := &pb.IDInput{
		ID: input.ID,
	}
	getSgOutput, err := h.poolHandler.Sg().Get(getSgInput, ctx)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case StorageRecordNotFoundErrCode:
				err = tkErr.New(cnt.OpskResourceRecordNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.OpskResource, "h.poolHandler.Sg().Get(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getSgInput),
		).Error(err.Error())
		return
	}

	output = &common.GetSecurityGroupOutput{
		ID:        getSgOutput.ID,
		Name:      getSgOutput.Name,
		UserID:    getSgOutput.User.ID,
		ProjectID: getSgOutput.Project.ID,
		Namespace: getSgOutput.Namespace,
	}
	return output, nil
}
