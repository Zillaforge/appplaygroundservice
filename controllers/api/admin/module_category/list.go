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
	// ListModuleCategoriesInput ...
	ListModuleCategoriesInput struct {
		adminCom.Pagination
		Where []string `json:"where" form:"where"`
		_     struct{}
	}
	// ListModuleCategoriesOutput ...
	ListModuleCategoriesOutput struct {
		ModuleCategories []adminCom.ModuleCategoryInterface `json:"moduleCategories"`
		Total            int                                `json:"total"`
		_                struct{}
	}
)

// ListModuleCategories ...
func ListModuleCategories(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input  = &ListModuleCategoriesInput{}
		output = &ListModuleCategoriesOutput{
			ModuleCategories: []adminCom.ModuleCategoryInterface{},
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

	// call grpc ListModuleJoinModuleAcls
	listModuleJoinModuleAclsInput := &pb.ListModuleJoinModuleAclsInput{
		Limit:  int32(input.Limit),
		Offset: int32(input.Offset),
		Where:  input.Where,
	}
	listModuleJoinModuleAclsOutput, listModuleJoinModuleAclsErr := aps.ListModuleJoinModuleAcls(listModuleJoinModuleAclsInput, c)
	if listModuleJoinModuleAclsErr != nil {
		if e, ok := tkErr.IsError(listModuleJoinModuleAclsErr); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.ControllerWhereQueryInvalidErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.ListModuleJoinModuleAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModuleJoinModuleAclsInput),
		).Error(listModuleJoinModuleAclsErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// reduce by ModuleCatogory-Module Pair
	existPair := map[string]map[string]bool{} // {<module category id>: <module ids>}
	for _, data := range listModuleJoinModuleAclsOutput.Data {
		// add ModuleCategoryID and detail when first occur
		if _, ok := existPair[data.ModuleCategoryID]; !ok {
			a := &adminCom.ModuleCategory{}
			a.ExtractByViewProto(c, data)
			output.ModuleCategories = append(output.ModuleCategories, a)

			existPair[data.ModuleCategoryID] = make(map[string]bool)
		}
		// skip if no Module
		if data.ModuleID == "" {
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
		moduleCategory, ok := data.(*adminCom.ModuleCategory)
		if !ok {
			err = tkErr.New(cnt.AdminAPIInternalServerErr)
			zap.L().With(
				zap.String(cnt.Controller, "data.(*adminCom.ModuleCategory); !ok"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", data),
			).Error(listModuleJoinModuleAclsErr.Error())
			statusCode = http.StatusInternalServerError
			utility.ResponseWithType(c, statusCode, err)
			return
		}
		count := len(existPair[moduleCategory.ID])
		output.ModuleCategories[i].(*adminCom.ModuleCategory).ModuleCount = &count
	}

	output.Total = len(existPair)
	utility.ResponseWithType(c, statusCode, output)
}
