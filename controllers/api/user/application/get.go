package application

import (
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
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
