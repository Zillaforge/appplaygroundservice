package module

import (
	"AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
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
	CreateModuleInput struct {
		Name             string `json:"name" binding:"required"`
		Description      string `json:"description"`
		Public           bool   `json:"public"`
		ModuleCategoryID string `json:"moduleCategoryId" binding:"required"`
		CreatorID        string `json:"creatorId" binding:"required"`
	}

	CreateModuleOutput adminCom.ModuleInterface
)

// CreateModule ...
func CreateModule(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input  = &CreateModuleInput{}
		output = (CreateModuleOutput)(&adminCom.Module{})
		err    error

		statusCode int = http.StatusCreated
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

	// check if user existed
	getUserInput := &authCom.GetUserInput{ID: input.CreatorID, Cacheable: true}
	_, err = authentication.Use().GetUser(c, getUserInput)
	if err != nil {
		if _err, ok := tkErr.IsError(err); ok {
			switch _err.Code() {
			case cnt.AuthUserNotFoundErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIUserNotFoundErr, input.CreatorID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "authentication.Use().GetUser(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getUserInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	createModuleInput := &pb.ModuleInfo{
		Name:             input.Name,
		Description:      input.Description,
		Public:           input.Public,
		ModuleCategoryID: input.ModuleCategoryID,
		CreatorID:        input.CreatorID,
	}
	createModuleOutput, createModuleErr := aps.CreateModule(createModuleInput, c)
	if createModuleErr != nil {
		if e, ok := tkErr.IsError(createModuleErr); ok {
			switch e.Code() {
			case cCnt.GRPCModuleExistErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIModuleAlreadyExistErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCModuleCategoryNotFoundErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIModuleCategoryNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.CreateModule(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createModuleInput),
		).Error(createModuleErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, createModuleOutput)
	utility.ResponseWithType(c, statusCode, output)
}
