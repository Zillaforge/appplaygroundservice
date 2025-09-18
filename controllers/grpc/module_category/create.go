package modulecategory

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) CreateModuleCategory(ctx context.Context, input *pb.ModuleCategoryInfo) (output *pb.ModuleCategoryInfo, err error) {
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

	createInput := &storCom.CreateModuleCategoryInput{
		ModuleCategory: tables.ModuleCategory{
			ID:          input.ID,
			Name:        input.Name,
			Description: input.Description,
			CreatorID:   input.CreatorID,
		},
	}

	createOutput, err := storages.Use().CreateModuleCategory(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageModuleCategoryExistErrCode:
				err = tkErr.New(cCnt.GRPCModuleCategoryExistErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateModuleCategory()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&createOutput.ModuleCategory)
	return
}
