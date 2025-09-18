package modulecategory

import (
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
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

type GetModuleCategoryInput struct {
	ID string `json:"-"`
	_  struct{}
}

type GetModuleCategoryOutput struct {
	userCom.ModuleCategory
	_ struct{}
}

func GetModuleCategory(c *gin.Context) {
	var (
		input      = &GetModuleCategoryInput{ID: c.GetString(cnt.CtxModuleCategoryID)}
		output     = &GetModuleCategoryOutput{}
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

	getModuleCategoryInput := &pb.GetInput{
		ID: input.ID,
	}
	getModuleCategoryOutput, err := aps.GetModuleCategory(getModuleCategoryInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetModuleCategory(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getModuleCategoryInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ModuleCategory.ExtractByProto(c, getModuleCategoryOutput)
	utility.ResponseWithType(c, statusCode, output)
}
