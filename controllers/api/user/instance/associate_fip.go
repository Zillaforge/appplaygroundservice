package instance

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
	userCom "AppPlaygroundService/controllers/api/user/common"
	"AppPlaygroundService/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type AssociateFloatingIPInput struct {
	ID           string `json:"-"`
	FloatingIPID string `json:"floatingIpId" binding:"required"`
	_            struct{}
}

type AssociateFloatingIPOutput struct {
	userCom.Instance
	_ struct{}
}

func AssociateFloatingIP(c *gin.Context) {
	var (
		input      = &AssociateFloatingIPInput{ID: c.GetString(cnt.CtxInstanceID)}
		output     = &AssociateFloatingIPOutput{}
		err        error
		requestID      = utility.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if err = c.ShouldBindWith(input, binding.JSON); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindWith()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	associateFloatingIPInput := &pb.UpdateFIPInput{
		ID:           input.ID,
		FloatingIPID: input.FloatingIPID,
	}
	associateFloatingIPOutput, err := aps.AssociateFloatingIP(associateFloatingIPInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.UserAPIInstanceNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCFloatingIPNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.UserAPIFloatingIPNotFoundErr, input.FloatingIPID)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCInstanceAlreadyHasFIPErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIInstanceAlreadyHasFIPErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCFIPCannotBeUseErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIFIPCannotBeUseErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.AssociateFloatingIP(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", associateFloatingIPInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(associateFloatingIPOutput)
	utility.ResponseWithType(c, statusCode, output)
}
