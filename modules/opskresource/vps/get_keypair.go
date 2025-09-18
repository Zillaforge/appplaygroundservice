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

func (h *Handler) GetKeypair(ctx context.Context, input *common.GetKeypairInput) (output *common.GetKeypairOutput, err error) {
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

	getKeypairInput := &pb.IDInput{
		ID: input.ID,
	}
	getKeypairOutput, err := h.poolHandler.Keypair().Get(getKeypairInput, ctx)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case StorageRecordNotFoundErrCode:
				err = tkErr.New(cnt.OpskResourceRecordNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.OpskResource, "h.poolHandler.Keypair().Get(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getKeypairInput),
		).Error(err.Error())
		return
	}

	output = &common.GetKeypairOutput{
		ID:     getKeypairOutput.ID,
		Name:   getKeypairOutput.Name,
		UserID: getKeypairOutput.User.ID,
	}
	return output, nil
}
