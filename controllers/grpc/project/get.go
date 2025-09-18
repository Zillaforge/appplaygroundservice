package project

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) GetProject(ctx context.Context, input *pb.GetInput) (output *pb.ProjectInfo, err error) {
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

	getInput := &storCom.GetProjectInput{
		ID: input.ID,
	}

	getOutput, err := storages.Use().GetProject(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageProjectNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCProjectNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetProject()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&getOutput.Project)
	return
}
