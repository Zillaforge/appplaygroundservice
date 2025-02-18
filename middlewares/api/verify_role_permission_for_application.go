package api

import (
	cnt "AppPlaygroundService/constants"
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

/*
Check if user's role can access application

errors:
- 12000000(internal server error)
- 12000007(application (%s) not found)
*/
func VerifyRolePermissionForApplication(c *gin.Context) {
	var (
		input        = &ResourceIDInput{ID: c.Param(cnt.ParamApplicationID)}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusOK
		supportRoles     = map[string]bool{
			cnt.TenantOwner.String(): true,
			cnt.TenantAdmin.String(): true,
		}
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"error":      &err,
		"statusCode": &statusCode,
	})

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

	// tenant-owner and tenant-admin allow to get all in project
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] {
		if getOutput.Application.CreatorID != c.GetString(cnt.CtxUserID) {
			statusCode = http.StatusNotFound
			err = tkErr.New(cnt.MidApplicationNotFoundErr, input.ID)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
	}

	c.Next()
}
