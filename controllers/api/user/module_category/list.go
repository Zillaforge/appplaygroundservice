package modulecategory

import (
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
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

type ListModuleCategoriesInput struct {
	userCom.Pagination

	ProjectID string `json:"-"`
	_         struct{}
}

type ListModuleCategoriesOutput struct {
	ModuleCategories []userCom.ModuleCategory `json:"moduleCategories"`
	Total            int                      `json:"total"`
	_                struct{}
}

func ListModuleCategories(c *gin.Context) {
	var (
		input = &ListModuleCategoriesInput{
			ProjectID: c.GetString(cnt.CtxProjectID),
		}
		output     = &ListModuleCategoriesOutput{}
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
		Limit:  int32(-1),
		Offset: int32(0),
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

	// reduce by ModuleCatogory-Module Pair
	existPair := map[string]map[string]bool{} // {<module category id>: <module ids>}
	output.ModuleCategories = []userCom.ModuleCategory{}
	for _, data := range listModuleJoinModuleAclsOutput.Data {
		// add ModuleCategoryID and detail when first occur
		if _, ok := existPair[data.ModuleCategoryID]; !ok {
			m := userCom.ModuleCategory{}
			m.ExtractByViewProto(c, data)
			output.ModuleCategories = append(output.ModuleCategories, m)

			existPair[data.ModuleCategoryID] = make(map[string]bool)
		}
		// skip if no Module or not permitted to see this Module.
		if data.ModuleID == "" || (!data.Public && data.AllowProjectID != input.ProjectID) {
			continue
		}
		// skip duplicated pair
		if _, ok := existPair[data.ModuleCategoryID][data.ModuleID]; ok {
			continue
		}
		existPair[data.ModuleCategoryID][data.ModuleID] = true
	}

	// Count Modules within each ModuleCategory
	for i, data := range output.ModuleCategories {
		count := len(existPair[data.ID])
		output.ModuleCategories[i].ModuleCount = &count
	}

	output.Total = len(existPair)
	utility.ResponseWithType(c, statusCode, output)
}
