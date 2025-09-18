package instance

import (
	cnt "AppPlaygroundService/constants"
	adminCom "AppPlaygroundService/controllers/api/admin/common"
	"AppPlaygroundService/utility"
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

type (
	// DisassociateFloatingIPInput ...
	DisassociateFloatingIPInput struct {
		ID string `json:"-"`
		_  struct{}
	}

	// DisassociateFloatingIPOutput ...
	DisassociateFloatingIPOutput adminCom.InstanceInterface
)

// DisassociateFloatingIP ...
func DisassociateFloatingIP(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &DisassociateFloatingIPInput{
			ID: c.GetString(cnt.CtxInstanceID),
		}
		output = (DisassociateFloatingIPOutput)(&adminCom.Instance{})
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

	disassociateFloatingIPInput := &pb.UpdateFIPInput{
		ID: input.ID,
	}
	disassociateFloatingIPOutput, disassociateFloatingIPErr := aps.DisassociateFloatingIP(disassociateFloatingIPInput, c)
	if disassociateFloatingIPErr != nil {
		if e, ok := tkErr.IsError(disassociateFloatingIPErr); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIInstanceNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCInstanceHasNoFIPErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIInstanceHasNoFIPErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.DisassociateFloatingIP(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", disassociateFloatingIPInput),
		).Error(disassociateFloatingIPErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(disassociateFloatingIPOutput)
	utility.ResponseWithType(c, statusCode, output)
}
