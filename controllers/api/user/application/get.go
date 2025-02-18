package application

import (
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
	"AppPlaygroundService/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type GetApplicationInput struct {
	ID string `json:"-"`
	_  struct{}
}

type GetApplicationOutput struct {
	userCom.Application
	_ struct{}
}

func GetApplication(c *gin.Context) {
	var (
		input      = &GetApplicationInput{ID: c.GetString(cnt.CtxApplicationID)}
		output     = &GetApplicationOutput{}
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

	getApplicationInput := &pb.GetWithLanguageInput{
		ID:       input.ID,
		Language: c.GetString(cnt.CtxLanguage),
	}
	getApplicationOutput, err := aps.GetApplication(getApplicationInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getApplicationInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Application.ExtractByProto(c, getApplicationOutput)
	utility.ResponseWithType(c, statusCode, output)
}
