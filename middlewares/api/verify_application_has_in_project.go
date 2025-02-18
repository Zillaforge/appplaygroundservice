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
Check application exists
endpoint: /project/:project-id/application/:application-id

errors:
- 12000000(internal server error)
- 12000007(application (%s) not found)
- 12000008(application (%s) is read only for users of other projects)
*/
func VerifyApplicationHasInProject(c *gin.Context) {
	var (
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamApplicationID)}
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
		err = tkErr.New(cnt.MidApplicationNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamApplicationID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// call GetApplication
	getInput := &pb.GetWithLanguageInput{
		ID: input.ID,
	}
	getOutput, err := aps.GetApplication(getInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCApplicationNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidApplicationNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "aps.GetApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// check permission
	if getOutput.Application.ProjectID != c.GetString(cnt.CtxProjectID) {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.MidApplicationIsReadOnlyErr, input.ID)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	c.Set(cnt.CtxApplicationID, input.ID)

	c.Next()
}
