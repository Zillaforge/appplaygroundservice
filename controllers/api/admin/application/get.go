package application

import (
	cnt "AppPlaygroundService/constants"
	adminCom "AppPlaygroundService/controllers/api/admin/common"
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
	// GetApplicationInput ...
	GetApplicationInput struct {
		ID string `json:"-"`
		_  struct{}
	}
	// GetApplicationOutput ...
	GetApplicationOutput adminCom.ApplicationInterface
)

// GetApplication ...
func GetApplication(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &GetApplicationInput{
			ID: c.GetString(cnt.CtxApplicationID),
		}
		output = (GetApplicationOutput)(&adminCom.Application{})
		err    error

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
	getApplicationOutput, getApplicationErr := aps.GetApplication(getApplicationInput, c)
	if getApplicationErr != nil {
		if e, ok := tkErr.IsError(getApplicationErr); ok {
			switch e.Code() {
			case cCnt.GRPCApplicationNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIApplicationNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getApplicationInput),
		).Error(getApplicationErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, getApplicationOutput)
	utility.ResponseWithType(c, statusCode, output)
}
