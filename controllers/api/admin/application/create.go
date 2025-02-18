package application

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
	CreateApplicationInput struct {
		Name        string                 `json:"name" binding:"required"`
		Description string                 `json:"description"`
		ModuleID    string                 `json:"moduleId" binding:"required"`
		Answers     map[string]interface{} `json:"answers" binding:"required"`
		Namespace   string                 `json:"namespace" binding:"required"`
		ProjectID   string                 `json:"projectId" binding:"required"`
		CreatorID   string                 `json:"creatorId" binding:"required"`
		Extra       map[string]interface{} `json:"extra"`
	}

	CreateApplicationOutput adminCom.ApplicationInterface
)

// CreateApplication ...
func CreateApplication(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input  = &CreateApplicationInput{}
		output = (CreateApplicationOutput)(&adminCom.Application{})
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

	if input.Answers == nil {
		input.Answers = map[string]interface{}{}
	}
	byteAnswers, marshalErr := json.Marshal(input.Answers)
	if marshalErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "json.Marshal(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input.Answers),
		).Error(marshalErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if input.Extra == nil {
		input.Extra = map[string]interface{}{}
	}
	byteExtra, marshalErr := json.Marshal(input.Extra)
	if marshalErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "json.Marshal(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input.Extra),
		).Error(marshalErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	createApplicationInput := &pb.CreateApplicationInput{
		Application: &pb.ApplicationInfo{
			Name:        input.Name,
			Description: input.Description,
			ModuleID:    input.ModuleID,
			Answers:     byteAnswers,
			Namespace:   input.Namespace,
			ProjectID:   input.ProjectID,
			CreatorID:   input.CreatorID,
			Language:    c.GetString(cnt.CtxLanguage),
			Extra:       byteExtra,
		},
	}
	createApplicationOutput, createApplicationErr := aps.CreateApplication(createApplicationInput, c)
	if createApplicationErr != nil {
		if e, ok := tkErr.IsError(createApplicationErr); ok {
			switch e.Code() {
			case cCnt.GRPCApplicationExistErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIApplicationAlreadyExistErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCModuleNotFoundErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIModuleNotFoundErr, input.ModuleID)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCProjectNotFoundErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIProjectNotFoundErr, input.ProjectID)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCQuizModuleErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIQuizModuleErr, e.Message())
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.CreateApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createApplicationInput),
		).Error(createApplicationErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErrWithInner, err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, createApplicationOutput)
	utility.ResponseWithType(c, statusCode, output)
}
