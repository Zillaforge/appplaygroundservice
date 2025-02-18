package module

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
	// ListModulesInput ...
	ListModulesInput struct {
		adminCom.Pagination
		Where    []string `json:"where" form:"where"`
		Language string   `json:"-"`
		_        struct{}
	}
	// ListModulesOutput ...
	ListModulesOutput struct {
		Modules []adminCom.ModuleInterface `json:"Modules"`
		Total   int                        `json:"total"`
		_       struct{}
	}
)

// ListModules ...
func ListModules(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &ListModulesInput{
			Language: c.GetString(cnt.CtxLanguage),
		}
		output = &ListModulesOutput{
			Modules: []adminCom.ModuleInterface{},
		}
		err error

		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if shouldBindQueryErr := c.ShouldBindQuery(input); shouldBindQueryErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindQuery(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(shouldBindQueryErr.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(shouldBindQueryErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// call grpc ListModules
	listModulesInput := &pb.ListWithLanguageInput{
		Limit:    int32(input.Limit),
		Offset:   int32(input.Offset),
		Where:    input.Where,
		Language: input.Language,
	}
	listModulesOutput, listModulesErr := aps.ListModules(listModulesInput, c)
	if listModulesErr != nil {
		if e, ok := tkErr.IsError(listModulesErr); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.ControllerWhereQueryInvalidErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.ListModules(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModulesInput),
		).Error(listModulesErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	for _, data := range listModulesOutput.Data {
		a := &adminCom.Module{}
		a.ExtractByProto(c, data)
		output.Modules = append(output.Modules, a)
	}
	output.Total = int(listModulesOutput.Count)

	utility.ResponseWithType(c, statusCode, output)

}
