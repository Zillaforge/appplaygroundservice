package module

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"github.com/google/uuid"

	"os"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

const ModuleActivatedState = "Activated"

func (m *Method) CreateModule(ctx context.Context, input *pb.ModuleInfo) (output *pb.ModuleDetail, err error) {
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

	// check the module category exist
	getModuleCategoryInput := &pb.GetInput{
		ID: input.ModuleCategoryID,
	}
	getModuleCategoryOutput, getModuleCategoryErr := aps.GetModuleCategory(getModuleCategoryInput, ctx)
	if getModuleCategoryErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "aps.GetModuleCategory()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getModuleCategoryInput),
		).Error(getModuleCategoryErr.Error())
		err = getModuleCategoryErr
		return
	}

	var id string
	// generate a new unique ID
	for {
		id = uuid.Must(uuid.NewRandom()).String()

		countInput := &storCom.CountModuleInput{ID: id}
		countOutput, countErr := storages.Use().CountModule(ctx, countInput)
		if countErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "storages.Use().CountModule()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", countInput),
			).Error(countErr.Error())
			err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(countErr)
			return
		}
		// check if id exists
		if countOutput.Count == 0 {
			break
		}
	}

	createInput := &storCom.CreateModuleInput{
		Module: tables.Module{
			ID:               id,
			Name:             input.Name,
			Description:      input.Description,
			ModuleCategoryID: input.ModuleCategoryID,
			Location:         utility.FormatShortModulePath(id),
			State:            ModuleActivatedState,
			Public:           input.Public,
			CreatorID:        input.CreatorID,
		},
	}

	createOutput, err := storages.Use().CreateModule(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageModuleExistErrCode:
				err = tkErr.New(cCnt.GRPCModuleExistErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateModule()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	// create module folder
	dirPath := utility.FormatModulePath(id)
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "os.MkdirAll"),
			zap.String(cnt.RequestID, requestID),
			zap.String("dirPath", dirPath),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&createOutput.Module)
	output.ModuleCategory = getModuleCategoryOutput

	return
}
