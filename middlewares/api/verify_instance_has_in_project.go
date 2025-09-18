package api

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

/*
Check instance exists
endpoint: /project/:project-id/instance/:instance-id

errors:
- 12000000(internal server error)
- 12000007(instance (%s) not found)
- 12000008(instance (%s) is read only for users of other projects)
*/
func VerifyInstanceHasInProject(c *gin.Context) {
	var (
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamInstanceID)}
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
		err = tkErr.New(cnt.MidInstanceNotFoundErr, input.ID)
		zap.L().With(
			zap.String(cnt.Middleware, "regexp.MatchString(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String(cnt.ParamInstanceID, input.ID),
		).Error(err.Error())
		statusCode = http.StatusNotFound
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// call GetInstance
	getInput := &pb.GetInput{
		ID: input.ID,
	}
	getOutput, err := aps.GetInstance(getInput, c)
	if err != nil {
		// Expected errors
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.MidInstanceNotFoundErr, input.ID)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		// Unexpected errors
		zap.L().With(
			zap.String(cnt.Middleware, "aps.GetInstance(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.MidInternalServerErrorErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// check permission
	if getOutput.Instance.ProjectID != c.GetString(cnt.CtxProjectID) {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.MidInstanceIsReadOnlyErr, input.ID)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	c.Set(cnt.CtxInstanceID, input.ID)

	c.Next()
}
