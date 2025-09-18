package application

import (
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
	"AppPlaygroundService/modules/fsmhandler"
	applicationCom "AppPlaygroundService/modules/fsmhandler/common/application"
	"AppPlaygroundService/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type RejectApplicationInput struct {
	ID string `json:"-"`
	_  struct{}
}

type RejectApplicationOutput struct {
	userCom.Application
	_ struct{}
}

func RejectApplication(c *gin.Context) {
	var (
		input        = &RejectApplicationInput{ID: c.GetString(cnt.CtxApplicationID)}
		output       = &RejectApplicationOutput{}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusOK
		supportRoles     = map[string]bool{
			cnt.TenantAdmin.String(): true,
		}
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	// Reject can only perform by TENANT_ADMIN
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] &&
		c.GetString(cnt.CtxCreator) != c.GetString(cnt.CtxUserID) {
		statusCode = http.StatusUnauthorized
		err = tkErr.New(cnt.UserAPIUnauthorizedOpErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	err = fsmhandler.Application.Event(c, input.ID, applicationCom.RejectEvent)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "fsmhandler.Application.Event(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.String("applicationID", input.ID),
			zap.String("event", applicationCom.RejectEvent),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// get application
	getApplicationInput := &pb.GetWithLanguageInput{
		ID:       input.ID,
		Language: c.GetString(cnt.CtxLanguage),
	}
	getApplicationOutput, err := aps.GetApplication(getApplicationInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getApplicationInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if lbmPublishErr := lbmPublish(c, getApplicationOutput, c.GetString(cnt.CtxUserID), false); lbmPublishErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "lbmPublish()"),
			zap.Any("app", getApplicationOutput),
			zap.String("reviewerID", c.GetString(cnt.CtxUserID)),
		).Error(lbmPublishErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, getApplicationOutput)
	utility.ResponseWithType(c, statusCode, output)
}
