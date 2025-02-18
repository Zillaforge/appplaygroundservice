package moduleacl

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (m *Method) CreateModuleAclBatch(ctx context.Context, input *pb.ModuleAclBatchInfo) (output *pb.ModuleAclBatchInfo, err error) {
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

	if len(input.Data) == 0 {
		return
	}
	createInput := &storCom.CreateModuleAclBatchInput{}
	for _, data := range input.Data {
		createInput.ModuleAcls = append(createInput.ModuleAcls, tables.ModuleAcl{
			ModuleID:  data.ModuleID,
			ProjectID: data.ProjectID,
		})
	}

	createOutput, err := storages.Use().CreateModuleAclBatch(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageModuleAclNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCModuleAclNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateModuleAclBatch()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = &pb.ModuleAclBatchInfo{
		Data: []*pb.ModuleAclInfo{},
	}
	for _, moduleAcl := range createOutput.ModuleAcls {
		output.Data = append(output.Data, m.storage2pb(&moduleAcl))
	}
	return
}
