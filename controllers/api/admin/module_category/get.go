package module_category

import (
	cnt "AppPlaygroundService/constants"
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
	// GetModuleCategoryInput ...
	GetModuleCategoryInput struct {
		ID string `json:"-"`
		_  struct{}
	}
	// GetModuleCategoryOutput ...
	GetModuleCategoryOutput adminCom.ModuleCategoryInterface
)

// GetModuleCategory ...
func GetModuleCategory(c *gin.Context) {

	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &GetModuleCategoryInput{
			ID: c.GetString(cnt.CtxModuleCategoryID),
		}
		output = (GetModuleCategoryOutput)(&adminCom.ModuleCategory{})
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

	getModuleCategoryInput := &pb.GetInput{
		ID: input.ID,
	}
	getModuleCategoryOutput, getModuleCategoryErr := aps.GetModuleCategory(getModuleCategoryInput, c)
	if getModuleCategoryErr != nil {
		if e, ok := tkErr.IsError(getModuleCategoryErr); ok {
			switch e.Code() {
			case cCnt.GRPCModuleCategoryNotFoundErrCode:
				statusCode = http.StatusNotFound
				err = tkErr.New(cnt.AdminAPIModuleCategoryNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetModuleCategory(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getModuleCategoryInput),
		).Error(getModuleCategoryErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, getModuleCategoryOutput)
	utility.ResponseWithType(c, statusCode, output)
}
