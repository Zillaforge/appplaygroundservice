package api

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type (
	VerifyApplicationInput struct {
		ID string
		_  struct{}
	}
)

func VerifyApplication(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &VerifyApplicationInput{}
		err   error

		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"error":      &err,
		"statusCode": &statusCode,
	})

	input.ID = c.Param(cnt.ParamApplicationID)

	getApplicationInput := &pb.GetWithLanguageInput{
		ID: input.ID,
	}
	if _, getApplicationErr := aps.GetApplication(getApplicationInput, c); getApplicationErr != nil {
		if e, ok := tkErr.IsError(getApplicationErr); ok {
			switch e.Code() {
			case cCnt.GRPCApplicationNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidApplicationNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
			zap.L().With(
				zap.String(cnt.Middleware, "aps.GetApplication()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getApplicationInput),
			).Error(getApplicationErr.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.MidInternalServerErrorErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
	}
	c.Set(cnt.CtxApplicationID, input.ID)
	c.Next()
}
