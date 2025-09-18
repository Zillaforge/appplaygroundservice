package module_category

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
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
	// UpdateModuleCategoryInput ...
	UpdateModuleCategoryInput struct {
		ID          string `json:"-"`
		Description string `json:"description"`
		_           struct{}
	}

	// UpdateFlavorOutput ...
	UpdateModuleCategoryOutput adminCom.ModuleCategoryInterface
)

// UpdateModuleCategory ...
func UpdateModuleCategory(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &UpdateModuleCategoryInput{
			ID: c.GetString(cnt.CtxModuleCategoryID),
		}
		output = (UpdateModuleCategoryOutput)(&adminCom.ModuleCategory{})
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

	updateModuleCategoryInput := &pb.UpdateModuleCategoryInput{
		ID:          input.ID,
		Description: &input.Description,
	}
	updateModuleCategoryOutput, updateModuleCategoryErr := aps.UpdateModuleCategory(updateModuleCategoryInput, c)
	if updateModuleCategoryErr != nil {
		if e, ok := tkErr.IsError(updateModuleCategoryErr); ok {
			switch e.Code() {
			case cCnt.GRPCModuleCategoryNotFoundErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIModuleCategoryNotFoundErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.UpdateModuleCategory(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateModuleCategoryInput),
		).Error(updateModuleCategoryErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, updateModuleCategoryOutput)
	utility.ResponseWithType(c, statusCode, output)

}
