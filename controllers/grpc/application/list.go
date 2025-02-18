package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/grpc"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/querydecoder"
	"context"

	"go.uber.org/zap"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (m *Method) ListApplications(ctx context.Context, input *pb.ListWithLanguageInput) (output *pb.ListApplicationsOutput, err error) {
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

	whereInput := storCom.ListApplicationsWhere{}
	if err := querydecoder.ShouldBindWhereSlice(&whereInput, input.Where); err != nil {
		if e, ok := tkErr.IsError(grpc.WhereErrorParser(err)); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				return output, e
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "querydecoder.ShouldBindWhereSlice(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input.where", input.Where),
		).Error(err.Error())
		return &pb.ListApplicationsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	listApplicationsInput := &storCom.ListApplicationsInput{
		Pagination: storCom.Paginate(input.Limit, input.Offset),
		Where:      whereInput,
	}
	listApplicationsOutput, err := storages.Use().ListApplications(ctx, listApplicationsInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().ListApplications(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listApplicationsInput),
		).Error(err.Error())
		return &pb.ListApplicationsOutput{}, tkErr.New(cCnt.GRPCInternalServerErr)
	}

	output = &pb.ListApplicationsOutput{
		Count: listApplicationsOutput.Count,
	}

	for _, application := range listApplicationsOutput.Applications {
		if input.Language != "" {
			application.Answers, err = m.getAnswers(
				string(application.Answers),
				application.ModuleID,
				input.Language)
			if err != nil {
				zap.L().With(
					zap.String(cnt.GRPC, "m.getAnswers()"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("ansStr", string(application.Answers)),
					zap.String("language", input.Language),
				).Error(err.Error())
				return
			}
		}
		output.Data = append(output.Data, m.storage2pb(&application))
	}

	return
}
