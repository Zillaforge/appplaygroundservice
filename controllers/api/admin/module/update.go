package module

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
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
	// UpdateModuleInput ...
	UpdateModuleInput struct {
		ID          string  `json:"-"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Public      *bool   `json:"public"`
		State       *string `json:"state"`
		Language    string  `json:"-"`
		_           struct{}
	}

	// UpdateModuleOutput ...
	UpdateModuleOutput adminCom.ModuleInterface
)

// UpdateModule ...
func UpdateModule(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &UpdateModuleInput{
			ID:       c.GetString(cnt.CtxModuleID),
			Language: c.GetString(cnt.CtxLanguage),
		}
		output = (UpdateModuleOutput)(&adminCom.Module{})
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

	updateModuleInput := &pb.UpdateModuleInput{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		State:       input.State,
		Public:      input.Public,
		Language:    input.Language,
	}
	updateModuleOutput, updateModuleErr := aps.UpdateModule(updateModuleInput, c)
	if updateModuleErr != nil {
		if e, ok := tkErr.IsError(updateModuleErr); ok {
			switch e.Code() {
			case cCnt.GRPCModuleExistErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIModuleAlreadyExistErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.UpdateModule(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateModuleInput),
		).Error(updateModuleErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, updateModuleOutput)
	utility.ResponseWithType(c, statusCode, output)
}
