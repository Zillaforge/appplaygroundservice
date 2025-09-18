package module

import (
	cnt "AppPlaygroundService/constants"
	adminCom "AppPlaygroundService/controllers/api/admin/common"
	"AppPlaygroundService/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type (
	// GetModuleInput ...
	GetModuleInput struct {
		ID       string `json:"-"`
		Language string `json:"-"`
		_        struct{}
	}
	// GetModuleOutput ...
	GetModuleOutput adminCom.ModuleInterface
)

// GetModule ...
func GetModule(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &GetModuleInput{
			ID:       c.GetString(cnt.CtxModuleID),
			Language: c.GetString(cnt.CtxLanguage),
		}
		output = (GetModuleOutput)(&adminCom.Module{})
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

	getModuleInput := &pb.GetWithLanguageInput{
		ID:       input.ID,
		Language: input.Language,
	}
	getModuleOutput, getModuleErr := aps.GetModule(getModuleInput, c)
	if getModuleErr != nil {
		if e, ok := tkErr.IsError(getModuleErr); ok {
			switch e.Code() {
			case cCnt.GRPCModuleNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIModuleNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCGetModuleQuestionsFailedErrCode:
				statusCode = http.StatusInternalServerError
				err = tkErr.New(cnt.AdminAPIGetModuleQuestionsFailedErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetModule(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getModuleInput),
		).Error(getModuleErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, getModuleOutput)
	utility.ResponseWithType(c, statusCode, output)
}
