package module_acl

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
	adminCom "AppPlaygroundService/controllers/api/admin/common"
	"AppPlaygroundService/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type (
	// UpdateModuleACLInput ...
	UpdateModuleACLInput struct {
		ModuleID   string   `json:"-"`
		ProjectIDs []string `json:"projectIds" binding:"required"`
		_          struct{}
	}

	// UpdateModuleACLOutput ...
	UpdateModuleACLOutput adminCom.ModuleACLInterface
)

// UpdateModuleACL ...
func UpdateModuleACL(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input = &UpdateModuleACLInput{
			ModuleID: c.GetString(cnt.CtxModuleID),
		}
		output = (UpdateModuleACLOutput)(&adminCom.ModuleACL{})
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

	// get the all projects
	listProjectsInput := &pb.ListInput{}
	listProjectsOutput, listProjectsErr := aps.ListProjects(listProjectsInput, c)
	if listProjectsErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.ListProjects(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listProjectsInput),
		).Error(listProjectsErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	var projectLists []string
	for _, projectInfo := range listProjectsOutput.Data {
		projectLists = append(projectLists, projectInfo.ID)
	}
	// check the project exist
	for _, project := range input.ProjectIDs {
		if !slices.Contains(projectLists, project) {
			statusCode = http.StatusBadRequest
			err = tkErr.New(cnt.AdminAPIProjectNotFoundErr, project)
			utility.ResponseWithType(c, statusCode, err)
			return
		}
	}

	// delete the ModuleAcl
	deleteModuleAclInput := &pb.DeleteInput{
		Where: []string{fmt.Sprintf("module-id=%s", input.ModuleID)},
	}
	if _, deleteModuleAclErr := aps.DeleteModuleAcl(deleteModuleAclInput, c); deleteModuleAclErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.DeleteModuleAcl(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", deleteModuleAclInput),
		).Error(deleteModuleAclErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// create the ModuleAcl in batch
	moduleAclInfos := []*pb.ModuleAclInfo{}
	for _, project := range input.ProjectIDs {
		moduleAclInfo := &pb.ModuleAclInfo{
			ModuleID:  input.ModuleID,
			ProjectID: project,
		}
		moduleAclInfos = append(moduleAclInfos, moduleAclInfo)
	}
	createModuleAclBatchInput := &pb.ModuleAclBatchInfo{
		Data: moduleAclInfos,
	}
	createModuleAclBatchOutput, createModuleAclBatchErr := aps.CreateModuleAclBatch(createModuleAclBatchInput, c)
	if createModuleAclBatchErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.CreateModuleAclBatch(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createModuleAclBatchInput),
		).Error(createModuleAclBatchErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}
	if len(createModuleAclBatchOutput.Data) == 0 {
		emptyData := &pb.ModuleAclInfo{
			ModuleID: input.ModuleID,
		}
		createModuleAclBatchOutput.Data = append(createModuleAclBatchOutput.Data, emptyData)
	}
	output.ExtractByProto(createModuleAclBatchOutput.Data)
	utility.ResponseWithType(c, statusCode, output)
}
