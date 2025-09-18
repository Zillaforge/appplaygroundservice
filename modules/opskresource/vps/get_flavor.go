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

func (h *Handler) GetFlavor(ctx context.Context, input *common.GetFlavorInput) (output *common.GetFlavorOutput, err error) {
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

	getFlavorInput := &pb.IDInput{
		ID: input.ID,
	}
	getFlavorOutput, err := h.poolHandler.Flavor().Get(getFlavorInput, ctx)
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
			zap.Any("input", getFlavorInput),
		).Error(err.Error())
		return
	}

	// 預設 GPU 資訊為空值，避免 nil panic
	var gpuInfo common.GpuInfo
	if getFlavorOutput.GPU != nil {
		gpuInfo = common.GpuInfo{
			Model:  getFlavorOutput.GPU.Model,
			Count:  getFlavorOutput.GPU.Count,
			IsVgpu: getFlavorOutput.GPU.IsVGPU,
		}
	}

	output = &common.GetFlavorOutput{
		ID:     getFlavorOutput.ID,
		Name:   getFlavorOutput.Name,
		AZ:     getFlavorOutput.AZ,
		Public: getFlavorOutput.Public,
		Vcpu:   getFlavorOutput.VCPU,
		Memory: getFlavorOutput.Memory,
		Disk:   getFlavorOutput.Disk,
		Gpu:    gpuInfo,
	}
	return output, nil
}
