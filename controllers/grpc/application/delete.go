package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/grpc"
	"AppPlaygroundService/modules/application"
	appCom "AppPlaygroundService/modules/application/common"
	fsmCom "AppPlaygroundService/modules/fsmhandler/common/application"
	"AppPlaygroundService/modules/opskresource"
	opskCom "AppPlaygroundService/modules/opskresource/common"
	"AppPlaygroundService/modules/opskresource/vps"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/querydecoder"
	"context"
	"time"

	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) DeleteApplication(ctx context.Context, input *pb.DeleteApplicationInput) (output *pb.DeleteOutput, err error) {
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

	// binding the where parameter
	whereInput := storCom.DeleteApplicationWhere{}
	if err = querydecoder.ShouldBindWhereSlice(&whereInput, input.Where); err != nil {
		if e, ok := tkErr.IsError(grpc.WhereErrorParser(err)); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				return output, e
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "querydecoder.ShouldBindWhereSlice(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input.Where", input.Where),
		).Error(err.Error())
		return output, tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
	}

	// get application
	getInput := &storCom.GetApplicationInput{
		ID: *whereInput.ID,
	}
	getOutput, err := storages.Use().GetApplication(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageApplicationNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCApplicationNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetApplication()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	if getOutput.Application.State == fsmCom.ProcessState {
		err = tkErr.New(cCnt.GRPCApplicationIsProcessingErr)
		zap.L().With(
			zap.String(cnt.GRPC, "getOutput.Application.State == fsmCom.ProcessState"),
			zap.String(cnt.RequestID, requestID),
		).Error(err.Error())
		return
	}

	// check and disassociate vps fip
	listInstancesInput := &storCom.ListInstancesInput{
		Pagination: storCom.Paginate(int32(-1), int32(0)),
		Where: storCom.ListInstancesWhere{
			ApplicationID: whereInput.ID,
		},
	}
	listInstancesOutput, err := storages.Use().ListInstances(ctx, listInstancesInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListInstances(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInstancesInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr)
		return
	}

	for _, instance := range listInstancesOutput.Instances {
		if instance.FloatingIPID != "" {
			// remove disassociated floating IP from VPS
			updateFIPInput := &opskCom.UpdateFloatingIPStatusInput{
				Action:       vps.ActionDisassociate,
				FloatingIPID: instance.FloatingIPID,
				IAMAuth: opskCom.IAMAuthInfo{
					UserID:    instance.Application.CreatorID,
					ProjectID: instance.Application.ProjectID,
				},
			}
			if _, err = opskresource.Use().UpdateFloatingIPStatus(ctx, updateFIPInput); err != nil {
				zap.L().With(
					zap.String(cnt.GRPC, "opskresource.Use().UpdateFloatingIPStatus(...)"),
					zap.Any("input", updateFIPInput),
				).Error(err.Error())
				err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
				return
			}
		}
	}

	// destroy instance
	destroyInput := &appCom.DestroyInput{
		ModuleID:      getOutput.Application.ModuleID,
		ApplicationID: getOutput.Application.ID,
		ProjectID:     getOutput.Application.ProjectID,
	}

	if input.AsyncDestroy {
		go application.Use().Destroy(ctx, *destroyInput)
	} else {
		destroyErr := application.Use().Destroy(ctx, *destroyInput)
		if destroyErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "application.Use().Destroy()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", *destroyInput),
			).Error(destroyErr.Error())
		}
	}

	// delete application
	deleteInput := &storCom.DeleteApplicationInput{
		Where: whereInput,
	}

	deleteOutput, err := storages.Use().DeleteApplication(ctx, deleteInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageApplicationInUseErrCode:
				err = tkErr.New(cCnt.GRPCApplicationInUseErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().DeleteApplication()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	// 更新 Application 的 EndedAt
	endedAt := time.Now().Format(time.RFC3339)
	updateMeteringInput := &pb.UpdateMeteringInput{
		ApplicationID: *whereInput.ID,
		EndedAt:       &endedAt,
	}
	if _, updateMeteringErr := aps.UpdateMetering(updateMeteringInput, ctx); updateMeteringErr != nil {
		if e, ok := tkErr.IsError(updateMeteringErr); ok {
			switch e.Code() {
			case cCnt.GRPCApplicationNotFoundErrCode:
				return &pb.DeleteOutput{
					ID: deleteOutput.ID,
				}, nil
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "aps.UpdateMetering(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateMeteringInput),
		).Error(updateMeteringErr.Error())
		return nil, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.DeleteOutput{
		ID: deleteOutput.ID,
	}
	return
}
