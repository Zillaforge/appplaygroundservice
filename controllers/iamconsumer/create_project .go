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
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	iamPB "github.com/Zillaforge/pegasusiamclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func UnmarshalCreateProject(input *ecCom.Data) (output *iamPB.ProjectID) {
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

func CreateProject(ctx context.Context, input *iamPB.ProjectID) {
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

	createInput := &pb.ProjectInfo{
		ID: input.ID,
	}
	if _, err = aps.CreateProject(createInput, ctx); err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCProjectExistErr.Code():
				err = tkErr.New(cnt.TaskProjectExistErr, input.ID)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.EventConsume, "aps.CreateProject(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		return
	}
}
