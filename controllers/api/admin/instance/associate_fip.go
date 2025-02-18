package instance

import (
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
	// AssociateFloatingIPInput ...
	AssociateFloatingIPInput struct {
		ID           string `json:"-"`
		FloatingIPID string `json:"floatingIpId" binding:"required"`
		_            struct{}
	}

	// AssociateFloatingIPOutput ...
	AssociateFloatingIPOutput adminCom.InstanceInterface
)

// AssociateFloatingIP ...
func AssociateFloatingIP(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &AssociateFloatingIPInput{
			ID: c.GetString(cnt.CtxInstanceID),
		}
		output = (AssociateFloatingIPOutput)(&adminCom.Instance{})
		err    error

		statusCode int = http.StatusOK
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

	associateFloatingIPInput := &pb.UpdateFIPInput{
		ID:           input.ID,
		FloatingIPID: input.FloatingIPID,
	}
	associateFloatingIPOutput, associateFloatingIPErr := aps.AssociateFloatingIP(associateFloatingIPInput, c)
	if associateFloatingIPErr != nil {
		if e, ok := tkErr.IsError(associateFloatingIPErr); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIInstanceNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCFloatingIPNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIFloatingIPNotFoundErr, input.FloatingIPID)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCInstanceAlreadyHasFIPErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIInstanceAlreadyHasFIPErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCFIPCannotBeUseErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIFIPCannotBeUseErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.AssociateFloatingIP(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", associateFloatingIPInput),
		).Error(associateFloatingIPErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(associateFloatingIPOutput)
	utility.ResponseWithType(c, statusCode, output)
}
