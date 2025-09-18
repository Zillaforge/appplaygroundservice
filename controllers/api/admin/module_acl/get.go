package module_acl

import (
	cnt "AppPlaygroundService/constants"
	adminCom "AppPlaygroundService/controllers/api/admin/common"
	"AppPlaygroundService/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type (
	// GetModuleACLInput ...
	GetModuleACLInput struct {
		ModuleID string `json:"-"`
		_        struct{}
	}
	// GetModuleACLOutput ...
	GetModuleACLOutput adminCom.ModuleACLInterface
)

// GetModuleACL ...
func GetModuleACL(c *gin.Context) {

	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &GetModuleACLInput{
			ModuleID: c.GetString(cnt.CtxModuleID),
		}
		output = (GetModuleACLOutput)(&adminCom.ModuleACL{})
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

	listModuleAclsInput := &pb.ListInput{
		Where: []string{fmt.Sprintf("module-id=%s", input.ModuleID)},
	}
	listModuleAclsOutput, listModuleAclsErr := aps.ListModuleAcls(listModuleAclsInput, c)
	if listModuleAclsErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.ListModuleAcls(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listModuleAclsInput),
		).Error(listModuleAclsErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if listModuleAclsOutput.Count == 0 {
		emptyData := &pb.ModuleAclInfo{
			ModuleID: input.ModuleID,
		}
		listModuleAclsOutput.Data = append(listModuleAclsOutput.Data, emptyData)
	}
	output.ExtractByProto(listModuleAclsOutput.Data)
	utility.ResponseWithType(c, statusCode, output)
}
