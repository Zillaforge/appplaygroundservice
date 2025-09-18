package app_credential

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/grpc"
	"AppPlaygroundService/modules/openstack"
	opstkCom "AppPlaygroundService/modules/openstack/common"
	"AppPlaygroundService/modules/opstkidentity"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) DeleteAppCredential(ctx context.Context, input *pb.DeleteInput) (output *pb.DeleteOutput, err error) {
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

	whereInput := storCom.AppCredentialWhere{}
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

	// List app credentials
	listInput := &storCom.ListAppCredentialsInput{
		Where: whereInput,
	}
	listOutput, err := storages.Use().ListAppCredentials(ctx, listInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListAppCredentials(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	// Delete app credentials from OpenStack
	for _, appCred := range listOutput.AppCredentials {
		opstkUID, uidErr := opstkidentity.Use().GetOpstkUID(ctx, appCred.UserID)
		if uidErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "o.GetOpstkUID"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("userID", appCred.UserID),
			).Error(uidErr.Error())
			return
		}

		opstkInput := &opstkCom.DeleteAppCredentialInput{
			ID:          appCred.ID,
			OpstkUserID: opstkUID,
		}
		err = openstack.Namespace(appCred.Namespace).Admin().Keystone().DeleteAppCredential(ctx, opstkInput)
		if err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "openstack.Namespace(...).Keystone(...).DeleteAppCredential"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", opstkInput),
			).Error(err.Error())
		}
	}

	// Delete app credentials from storage
	deleteInput := &storCom.DeleteAppCredentialInput{
		Where: whereInput,
	}
	deleteOutput, err := storages.Use().DeleteAppCredential(ctx, deleteInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageAppCredentialNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCAppCredentialNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().DeleteAppCredential()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = &pb.DeleteOutput{
		ID: deleteOutput.ID,
	}
	return
}
