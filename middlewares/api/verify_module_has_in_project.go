package api

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

/*
Check module exists
endpoint: /project/:project-id/module/:module-id

errors:
- 12000000(internal server error)
- 12000007(module (%s) not found)
- 12000008(module (%s) is read only for users of other projects)
*/
func VerifyModuleHasInProject(c *gin.Context) {
	var (
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamModuleID)}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"error":      &err,
		"statusCode": &statusCode,
	})

	// validate input
	match, _ := regexp.MatchString(uuidRegexpString, input.ID)
	if !match {
		err = tkErr.New(cnt.MidModuleNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamModuleID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// call GetModule
	getInput := &pb.GetWithLanguageInput{
		ID: input.ID,
	}
	getOutput, err := aps.GetModule(getInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCModuleNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidModuleNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "aps.GetModule(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// check permission
	if !getOutput.Module.Public {
		listInput := &pb.ListInput{
			Where: []string{
				"module-id=" + input.ID,
				"project-id=" + c.GetString(cnt.CtxProjectID),
			},
		}
		listOutput, err := aps.ListModuleAcls(listInput, c)
		if err != nil {
			// Unexpected errors
			zap.L().With(
				zap.String(cnt.Middleware, "aps.ListModuleAcls(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", listInput),
			).Error(err.Error())
			statusCode = http.StatusInternalServerError
			err = tkErr.New(cnt.MidInternalServerErrorErr)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		if listOutput.Count == 0 {
			statusCode = http.StatusForbidden
			err = tkErr.New(cnt.MidModuleIsReadOnlyErr, input.ID)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
	}

	c.Set(cnt.CtxModuleID, input.ID)

	c.Next()
}
