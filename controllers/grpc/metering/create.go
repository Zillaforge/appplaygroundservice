package metering

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"
	"time"

	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (m *Method) CreateMetering(ctx context.Context, input *pb.MeteringInfo) (output *pb.MeteringInfo, err error) {
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

	createMeteringInput := &storCom.CreateMeteringInput{
		Metering: tables.Metering{
			ApplicationID: input.ApplicationID,
			Name:          input.Name,
			ProjectID:     input.ProjectID,
			Creator:       input.Creator,
			Instances:     input.Instances,
			CreatedAt: func(input string) (output time.Time) {
				parseOutput, parseErr := time.Parse(time.RFC3339, input)
				if parseErr != nil {
					return time.Now()
				}
				return parseOutput
			}(input.CreatedAt),
		},
	}
	createMeteringOutput, createMeteringErr := storages.Use().CreateMetering(ctx, createMeteringInput)
	if createMeteringErr != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateMetering(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createMeteringInput),
		).Error(createMeteringErr.Error())
		return
	}
	return m.storage2pb(&createMeteringOutput.Metering), nil
}
