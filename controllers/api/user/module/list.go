package module

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

type ListModulesInput struct {
	userCom.Pagination
	Where []string `json:"-"`

	ProjectID        string `json:"-"`
	ModuleCategoryID string `json:"-"`
	Language         string `json:"-"`
	_                struct{}
}

type ListModulesOutput struct {
	Modules []userCom.Module `json:"modules"`
	Total   int              `json:"total"`
	_       struct{}
}

func ListModules(c *gin.Context) {
	var (
		moduleCategoryID = c.GetString(cnt.CtxModuleCategoryID)
		input            = &ListModulesInput{
			Where: []string{
				"module-category-id==" + moduleCategoryID,
			},
			ProjectID:        c.GetString(cnt.CtxProjectID),
			ModuleCategoryID: moduleCategoryID,
			Language:         c.GetString(cnt.CtxLanguage),
		}
		output     = &ListModulesOutput{}
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

	listModuleJoinModuleAclsInput := &pb.ListModuleJoinModuleAclsInput{
		Limit:     int32(-1),
		Offset:    int32(0),
		Where:     input.Where,
		ProjectID: &input.ProjectID,
		Language:  input.Language,
	}
	listModuleJoinModuleAclsOutput, err := aps.ListModuleJoinModuleAcls(listModuleJoinModuleAclsInput, c)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				statusCode = http.StatusBadRequest
				if v, exist := e.Get("field"); exist {
					err = tkErr.New(cnt.UserAPIQueryNotSupportErr, "where", v)
				} else {
					err = tkErr.New(cnt.UserAPIIllegalWhereQueryFormatErr)
				}
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.ListModuleJoinModuleAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModuleJoinModuleAclsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Modules = []userCom.Module{}
	existID := map[string]bool{}
	for _, m := range listModuleJoinModuleAclsOutput.Data {
		if _, ok := existID[m.ModuleID]; !ok {
			existID[m.ModuleID] = true

			module := userCom.Module{}
			module.ExtractByViewProto(c, m)
			output.Modules = append(output.Modules, module)
		}
	}
	output.Total = len(existID)

	c.Writer.Header().Set(cnt.HdrContentLanguage, input.Language)
	utility.ResponseWithType(c, statusCode, output)
}
