package api

import (
	cnt "AppPlaygroundService/constants"
	util "AppPlaygroundService/utility"
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
Check module category exists

errors:
- 12000000(internal server error)
- 12000006(module category (%s) not found)
*/
func VerifyModuleCategory(c *gin.Context) {
	var (
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamModuleCategoryID)}
		err        error
		requestID      = util.MustGetContextRequestID(c)
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
		err = tkErr.New(cnt.MidModuleCategoryNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamModuleCategoryID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		util.ResponseWithType(c, statusCode, err)
		return
	}

	// call GetModuleCategory
	getInput := &pb.GetInput{
		ID: input.ID,
	}
	_, err = aps.GetModuleCategory(getInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCModuleCategoryNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidModuleCategoryNotFoundErr, input.ID)
				util.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "aps.GetModuleCategory(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		util.ResponseWithType(c, statusCode, err)
		return
	}

	c.Set(cnt.CtxModuleCategoryID, input.ID)
	c.Next()
}
