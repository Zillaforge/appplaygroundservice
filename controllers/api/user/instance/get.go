package instance

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

type GetInstanceInput struct {
	ID string `json:"-"`
	_  struct{}
}

type GetInstanceOutput struct {
	userCom.Instance
	_ struct{}
}

func GetInstance(c *gin.Context) {
	var (
		input      = &GetInstanceInput{ID: c.GetString(cnt.CtxInstanceID)}
		output     = &GetInstanceOutput{}
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

	getInstanceInput := &pb.GetInput{
		ID: input.ID,
	}
	getInstanceOutput, err := aps.GetInstance(getInstanceInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetInstance(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInstanceInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Instance.ExtractByProto(getInstanceOutput)
	utility.ResponseWithType(c, statusCode, output)
}
