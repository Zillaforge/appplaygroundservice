package instance

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
	adminCom "AppPlaygroundService/controllers/api/admin/common"
	"AppPlaygroundService/utility"
	"encoding/json"
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
	// UpdateInstanceInput ...
	UpdateInstanceInput struct {
		ID    string                  `json:"-"`
		Name  *string                 `json:"name"`
		Extra *map[string]interface{} `json:"extra"`
		_     struct{}
	}

	// UpdateInstanceOutput ...
	UpdateInstanceOutput adminCom.InstanceInterface
)

// UpdateInstance ...
func UpdateInstance(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &UpdateInstanceInput{
			ID: c.GetString(cnt.CtxInstanceID),
		}
		output = (UpdateInstanceOutput)(&adminCom.Instance{})
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

	if shouldBindJSONErr := c.ShouldBindJSON(input); shouldBindJSONErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindJSON(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(shouldBindJSONErr.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(shouldBindJSONErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	updateInstanceInput := &pb.UpdateInstanceInput{
		ID:   input.ID,
		Name: input.Name,
		Extra: func() []byte {
			if input.Extra != nil {
				extraByte, _ := json.Marshal(input.Extra)
				return extraByte
			}
			return nil
		}(),
	}
	updateInstanceOutput, updateInstanceErr := aps.UpdateInstance(updateInstanceInput, c)
	if updateInstanceErr != nil {
		if e, ok := tkErr.IsError(updateInstanceErr); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIInstanceNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.UpdateInstance(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInstanceInput),
		).Error(updateInstanceErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(updateInstanceOutput)
	utility.ResponseWithType(c, statusCode, output)
}
