package instance

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
	// GetModuleInput ...
	GetInstanceInput struct {
		ID string `json:"-"`
		_  struct{}
	}
	// GetModuleOutput ...
	GetInstanceOutput adminCom.InstanceInterface
)

// GetInstance ...
func GetInstance(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &GetInstanceInput{
			ID: c.GetString(cnt.CtxInstanceID),
		}
		output = (GetInstanceOutput)(&adminCom.Instance{})
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

	getInstanceInput := &pb.GetInput{
		ID: input.ID,
	}
	getInstanceOutput, getInstanceErr := aps.GetInstance(getInstanceInput, c)
	if getInstanceErr != nil {
		if e, ok := tkErr.IsError(getInstanceErr); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIInstanceNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetInstance(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInstanceInput),
		).Error(getInstanceErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(getInstanceOutput)
	utility.ResponseWithType(c, statusCode, output)
}
