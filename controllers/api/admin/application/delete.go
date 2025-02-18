package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"fmt"
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
	// DeleteApplicationInput ...
	DeleteApplicationInput struct {
		Where []string `json:"where" form:"where"`
		_     struct{}
	}
	// DeleteApplicationOutput ...
	DeleteApplicationOutput struct {
		_ struct{}
	}
)

// DeleteApplication ...
func DeleteApplication(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &DeleteApplicationInput{
			Where: []string{fmt.Sprintf("ID=%s", c.GetString(cnt.CtxApplicationID))},
		}
		output = &DeleteApplicationOutput{}
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

	deleteApplicationInput := &pb.DeleteApplicationInput{
		Where:        input.Where,
		AsyncDestroy: true,
	}
	if _, deleteApplicationErr := aps.DeleteApplication(deleteApplicationInput, c); deleteApplicationErr != nil {
		if e, ok := tkErr.IsError(deleteApplicationErr); ok {
			switch e.Code() {
			case cCnt.GRPCApplicationIsProcessingErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIApplicationIsProcessingErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.DeleteApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteApplicationInput),
		).Error(deleteApplicationErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, nil)
}
