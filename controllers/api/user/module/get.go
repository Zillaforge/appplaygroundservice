package module

import (
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
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

type GetModuleInput struct {
	ID       string `json:"-"`
	Language string `json:"-"`
	_        struct{}
}

type GetModuleOutput struct {
	userCom.Module
	_ struct{}
}

func GetModule(c *gin.Context) {
	var (
		input = &GetModuleInput{
			ID:       c.GetString(cnt.CtxModuleID),
			Language: c.GetString(cnt.CtxLanguage),
		}
		output     = &GetModuleOutput{}
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

	getModuleInput := &pb.GetWithLanguageInput{
		ID:       input.ID,
		Language: input.Language,
	}
	getModuleOutput, err := aps.GetModule(getModuleInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCModuleNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.UserAPIModuleNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCGetModuleQuestionsFailedErrCode:
				statusCode = http.StatusInternalServerError
				err = tkErr.New(cnt.UserAPIGetModuleQuestionsFailedErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetModule(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getModuleInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Module.ExtractByProto(c, getModuleOutput)

	c.Writer.Header().Set(cnt.HdrContentLanguage, input.Language)
	utility.ResponseWithType(c, statusCode, output)
}
