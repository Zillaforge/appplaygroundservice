package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
	userCom "AppPlaygroundService/controllers/api/user/common"
	"AppPlaygroundService/utility"
	"encoding/json"
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

type CreateApplicationInput struct {
	Name        string                 `json:"name" binding:"required"`
	ModuleID    string                 `json:"moduleId" binding:"required"`
	Answers     map[string]interface{} `json:"answers" binding:"required"`
	Description string                 `json:"description"`
	Extra       map[string]interface{} `json:"extra"`

	Namespace string `json:"-"`
	ProjectID string `json:"-"`
	CreatorID string `json:"-"`
	UpdaterID string `json:"-"`

	_ struct{}
}

type CreateApplicationOutput struct {
	userCom.Application
	_ struct{}
}

func CreateApplication(c *gin.Context) {
	var (
		input      = &CreateApplicationInput{}
		output     = &CreateApplicationOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusCreated
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if err = c.ShouldBindJSON(input); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindJSON(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
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
		err = tkErr.New(cnt.UserAPIInternalServerErr)
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
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	createApplicationInput := &pb.CreateApplicationInput{
		Application: &pb.ApplicationInfo{
			Name:        input.Name,
			ModuleID:    input.ModuleID,
			Answers:     byteAnswers,
			Description: input.Description,
			Namespace:   c.GetString(cnt.CtxNamespace),
			ProjectID:   c.GetString(cnt.CtxProjectID),
			CreatorID:   c.GetString(cnt.CtxUserID),
			UpdaterID:   c.GetString(cnt.CtxUserID),
			Language:    c.GetString(cnt.CtxLanguage),
			Extra:       byteExtra,
		},
		UnderReview: c.GetBool(cnt.CtxResourceReview),
	}
	createApplicationOutput, err := aps.CreateApplication(createApplicationInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCApplicationExistErr.Code():
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIApplicationExistErr, input.Name)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCModuleNotFoundErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIModuleNotFoundErr, input.ModuleID)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCProjectNotFoundErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIProjectNotFoundErr, input.ProjectID)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCQuizModuleErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIQuizModuleErr, e.Message())
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Controller, "aps.CreateApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createApplicationInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErrWithInner, err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Application.ExtractByProto(c, createApplicationOutput)
	utility.ResponseWithType(c, statusCode, output)
}
