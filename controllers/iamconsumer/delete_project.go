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
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	iamPB "github.com/Zillaforge/pegasusiamclient/pb"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func UnmarshalDeleteProject(input *ecCom.Data) (output *iamPB.ProjectID) {
	output = &iamPB.ProjectID{}
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

func DeleteProject(ctx context.Context, input *iamPB.ProjectID) {
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

	deleteInput := &pb.DeleteInput{
		Where: []string{"ID=" + input.ID},
	}
	if _, err = aps.DeleteProject(deleteInput, ctx); err != nil {
		zap.L().With(
			zap.String(cnt.EventConsume, "aps.DeleteProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteInput),
		).Error(err.Error())
		return
	}
}
