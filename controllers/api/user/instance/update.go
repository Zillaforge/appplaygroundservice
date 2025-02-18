package instance

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
	userCom "AppPlaygroundService/controllers/api/user/common"
	"AppPlaygroundService/utility"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type UpdateInstanceInput struct {
	ID    string                 `json:"-"`
	Extra map[string]interface{} `json:"extra" binding:"required"`
	_     struct{}
}

type UpdateInstanceOutput struct {
	userCom.Instance
	_ struct{}
}

func UpdateInstance(c *gin.Context) {
	var (
		input      = &UpdateInstanceInput{ID: c.GetString(cnt.CtxInstanceID)}
		output     = &UpdateInstanceOutput{}
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

	if err = c.ShouldBindWith(input, binding.JSON); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindWith()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	metaByte, marshalErr := json.Marshal(input.Extra)
	if marshalErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "json.Marshal(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input.Extra),
		).Error(marshalErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.ControllerInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	updateInput := &pb.UpdateInstanceInput{
		ID:    input.ID,
		Extra: metaByte,
	}
	updateInstanceOutput, err := aps.UpdateInstance(updateInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.UserAPIInstanceNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.UpdateInstance(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(updateInstanceOutput)
	utility.ResponseWithType(c, statusCode, output)
}
