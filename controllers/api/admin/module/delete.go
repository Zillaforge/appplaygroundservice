package module

import (
	"AppPlaygroundService/utility"
	"fmt"
	"net/http"

	cnt "AppPlaygroundService/constants"

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
	// DeleteModuleInput ...
	DeleteModuleInput struct {
		Where []string
		_     struct{}
	}
	// DeleteModuleOutput ...
	DeleteModuleOutput struct {
		_ struct{}
	}
)

// DeleteModule ...
func DeleteModule(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &DeleteModuleInput{
			Where: []string{fmt.Sprintf("ID=%s", c.GetString(cnt.CtxModuleID))},
		}
		output = &DeleteModuleOutput{}
		err    error

		statusCode int = http.StatusNoContent
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	deleteModuleInput := &pb.DeleteInput{
		Where: input.Where,
	}
	if _, deleteModuleErr := aps.DeleteModule(deleteModuleInput, c); deleteModuleErr != nil {
		if e, ok := tkErr.IsError(deleteModuleErr); ok {
			switch e.Code() {
			case cCnt.GRPCModuleInUseErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIModuleAlreadyInUseErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.DeleteModule(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteModuleInput),
		).Error(deleteModuleErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, nil)
}
