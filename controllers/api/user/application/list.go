package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
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

type ListApplicationsInput struct {
	userCom.Pagination
	Where []string `json:"where" form:"where"`

	ProjectID string `json:"-"`
	_         struct{}
}

type ListApplicationsOutput struct {
	Applications []userCom.Application `json:"applications"`
	Total        int                   `json:"total"`
	_            struct{}
}

func ListApplications(c *gin.Context) {
	var (
		input = &ListApplicationsInput{
			ProjectID: c.GetString(cnt.CtxProjectID),
		}
		output       = &ListApplicationsOutput{}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusOK
		supportRoles     = map[string]bool{
			cnt.TenantOwner.String(): true,
			cnt.TenantAdmin.String(): true,
		}
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	if err = c.ShouldBindQuery(input); err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "c.ShouldBindQuery(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("obj", input),
		).Error(err.Error())
		statusCode = http.StatusBadRequest
		err = api.Malformed(err)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// tenant-owner and tenant-admin allow to get all in project
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] {
		input.Where = append(input.Where, "creator-id="+c.GetString(cnt.CtxUserID))
	}

	input.Where = append(input.Where, "project-id="+input.ProjectID)
	input.Where = append(input.Where, "namespace="+c.GetString(cnt.CtxNamespace))

	listApplicationsInput := &pb.ListWithLanguageInput{
		Limit:    int32(input.Limit),
		Offset:   int32(input.Offset),
		Where:    input.Where,
		Language: c.GetString(cnt.CtxLanguage),
	}
	listApplicationsOutput, err := aps.ListApplications(listApplicationsInput, c)
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
			zap.String(cnt.Controller, "aps.ListApplications(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listApplicationsInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.Applications = []userCom.Application{}
	output.Total = int(listApplicationsOutput.Count)
	for _, data := range listApplicationsOutput.Data {
		m := userCom.Application{}
		m.ExtractByProto(c, data)
		output.Applications = append(output.Applications, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
