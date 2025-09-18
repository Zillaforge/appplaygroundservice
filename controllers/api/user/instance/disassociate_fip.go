package instance

import (
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
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

type DisassociateFloatingIPInput struct {
	ID string `json:"-"`
	_  struct{}
}

type DisassociateFloatingIPOutput struct {
	userCom.Instance
	_ struct{}
}

func DisassociateFloatingIP(c *gin.Context) {
	var (
		input      = &DisassociateFloatingIPInput{ID: c.GetString(cnt.CtxInstanceID)}
		output     = &DisassociateFloatingIPOutput{}
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

	disassociateFloatingIPInput := &pb.UpdateFIPInput{
		ID: input.ID,
	}
	disassociateFloatingIPOutput, err := aps.DisassociateFloatingIP(disassociateFloatingIPInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCInstanceNotFoundErr.Code():
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.UserAPIInstanceNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			case cCnt.GRPCInstanceHasNoFIPErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.UserAPIInstanceHasNoFIPErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.DisassociateFloatingIP(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", disassociateFloatingIPInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(disassociateFloatingIPOutput)
	utility.ResponseWithType(c, statusCode, output)
}
