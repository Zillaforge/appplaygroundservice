package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type GetLogsInput struct {
	ApplicationID string `json:"-"`
	_             struct{}
}

type GetLogsOutput struct {
	Logs string `json:"logs"`
	_    struct{}
}

func GetAppLogs(c *gin.Context) {
	var (
		input = &GetLogsInput{
			ApplicationID: c.GetString(cnt.CtxApplicationID),
		}
		output     = &GetLogsOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	getAppLogsInput := &pb.GetAppLogsInput{
		ApplicationID: input.ApplicationID,
	}
	getAppLogsOutput, err := aps.GetAppLogs(getAppLogsInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetLogs(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getAppLogsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Logs = getAppLogsOutput.Logs
	utility.ResponseWithType(c, statusCode, output)
}
