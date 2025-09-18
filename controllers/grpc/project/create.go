package project

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

/*
CreateProject calls storages interface to create a new project.
And returns the new project information.

errors:
- 14000000(internal server error)
- 14020000(project exist)
*/
func (m *Method) CreateProject(ctx context.Context, input *pb.ProjectInfo) (output *pb.ProjectInfo, err error) {
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

	createProjectInput := &storCom.CreateProjectInput{
		Project: tables.Project{
			ID: input.ID,
		},
	}
	createProjectOutput, err := storages.Use().CreateProject(ctx, createProjectInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageProjectExistErrCode:
				return &pb.ProjectInfo{}, tkErr.New(cCnt.GRPCProjectExistErr)
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createProjectInput),
		).Error(err.Error())
		return &pb.ProjectInfo{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	return m.storage2pb(&createProjectOutput.Project), nil
}
