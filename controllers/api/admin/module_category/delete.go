package module_category

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"fmt"
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
	// DeleteModuleCategoryInput ...
	DeleteModuleCategoryInput struct {
		Where []string `json:"where" form:"where"`
		_     struct{}
	}
	// DeleteModuleCategoryOutput ...
	DeleteModuleCategoryOutput struct {
		_ struct{}
	}
)

// DeleteModuleCategory ...
func DeleteModuleCategory(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &DeleteModuleCategoryInput{
			Where: []string{fmt.Sprintf("ID=%s", c.GetString(cnt.CtxModuleCategoryID))},
		}
		output = &DeleteModuleCategoryOutput{}
		err    error

		statusCode int = http.StatusNoContent
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	deleteModuleCategoryInput := &pb.DeleteInput{
		Where: input.Where,
	}
	if _, deleteModuleCategoryErr := aps.DeleteModuleCategory(deleteModuleCategoryInput, c); deleteModuleCategoryErr != nil {
		if e, ok := tkErr.IsError(deleteModuleCategoryErr); ok {
			switch e.Code() {
			case cCnt.GRPCModuleCategoryInUseErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.AdminAPIModuleCategoryAlreadyInUseErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.DeleteModuleCategory(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteModuleCategoryInput),
		).Error(deleteModuleCategoryErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	utility.ResponseWithType(c, statusCode, nil)

}
