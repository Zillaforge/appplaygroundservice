package project

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/grpc"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

/*
DeleteProject calls storages interface to delete a project.

errors:
- 14000000(internal server error)
*/

func (m *Method) DeleteProject(ctx context.Context, input *pb.DeleteInput) (output *pb.DeleteOutput, err error) {
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
	whereInput := storCom.DeleteProjectWhere{}
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

	// list application by project-id and delete them one-by-one.
	listApplicationsInput := &storCom.ListApplicationsInput{
		Pagination: storCom.Paginate(int32(-1), int32(0)),
		Where: storCom.ListApplicationsWhere{
			ProjectID: whereInput.ID,
		},
	}
	listApplicationsOutput, err := storages.Use().ListApplications(ctx, listApplicationsInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListApplications(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listApplicationsInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	for _, application := range listApplicationsOutput.Applications {
		deleteAppInput := &pb.DeleteApplicationInput{
			Where:        []string{"ID=" + application.ID},
			AsyncDestroy: false,
		}
		_, err = aps.DeleteApplication(deleteAppInput, ctx)
		if err != nil {
			zap.L().With(
				zap.String(cnt.Controller, "aps.DeleteApplication(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", deleteAppInput),
			).Error(err.Error())
			return
		}
	}

	// Delete app credentials by project-id
	deleteAppCredInput := &pb.DeleteInput{
		Where: []string{"project-id=" + *whereInput.ID},
	}
	if _, err = aps.DeleteAppCredential(deleteAppCredInput, ctx); err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "aps.DeleteAppCredential(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteAppCredInput),
		).Error(err.Error())
		return
	}

	// Delete project
	deleteInput := &storCom.DeleteProjectInput{
		Where: whereInput,
	}
	deleteOutput, err := storages.Use().DeleteProject(ctx, deleteInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageProjectNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCProjectNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().DeleteProject()"),
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
