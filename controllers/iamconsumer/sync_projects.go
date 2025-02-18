package iamconsumer

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/logger"
	ecCom "AppPlaygroundService/modules/eventconsume/common"
	"AppPlaygroundService/utility"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	iamPB "pegasus-cloud.com/aes/pegasusiamclient/pb"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
)

func UnmarshalSyncProjects(input *ecCom.Data) (output *iamPB.ListProjectOutput) {
	output = &iamPB.ListProjectOutput{}
	if input.Request != nil {
		switch v := input.Request.(type) {
		case string:
			decodedData, _ := base64.StdEncoding.DecodeString(v)
			json.Unmarshal(decodedData, output)
			logger.Use().Info(fmt.Sprintf("%s | %s | %s | %s",
				input.Metadata[tracer.RequestID],
				"FromIAM",
				input.Action,
				decodedData,
			))
		default:
			zap.L().Warn(fmt.Sprintf("Received the message format of %s action is invalid", input.Action))
		}
	}
	return output
}

func SyncProjects(ctx context.Context, input *iamPB.ListProjectOutput) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
		err       error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input": &input,
			"err":   &err,
		},
	)

	if err = aps.SyncProjects(ctx); err != nil {
		zap.L().With(
			zap.String(cnt.EventConsume, "aps.SyncProjects(...)"),
			zap.String(cnt.RequestID, requestID),
		).Error(err.Error())
		return
	}
}
