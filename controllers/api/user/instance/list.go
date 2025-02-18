package instance

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

type ListInstancesInput struct {
	userCom.Pagination
	Where []string `json:"where" form:"where"`

	ProjectID string `json:"-"`
	_         struct{}
}

type ListInstancesOutput struct {
	Instances []userCom.Instance `json:"instances"`
	Total     int                `json:"total"`
	_         struct{}
}

func ListInstances(c *gin.Context) {
	var (
		input        = &ListInstancesInput{}
		output       = &ListInstancesOutput{}
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

	listInstancesInput := &pb.ListInput{
		Limit:  int32(input.Limit),
		Offset: int32(input.Offset),
		Where:  append(input.Where, "project-id="+c.GetString(cnt.CtxProjectID)),
	}
	listInstancesOutput, err := aps.ListInstances(listInstancesInput, c)
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
			zap.String(cnt.Controller, "aps.ListInstances(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInstancesInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// tenant-owner and tenant-admin allow to get all in project
	supportedRole := false
	if role := c.GetString(cnt.CtxTenantRole); supportRoles[role] {
		supportedRole = true
	}
	userID := c.GetString(cnt.CtxUserID)

	output.Instances = []userCom.Instance{}
	output.Total = int(listInstancesOutput.Count)
	for _, data := range listInstancesOutput.Data {
		if !supportedRole {
			if data.Application.CreatorID != userID {
				continue
			}
		}
		m := userCom.Instance{}
		m.ExtractByProto(data)
		output.Instances = append(output.Instances, m)
	}
	utility.ResponseWithType(c, statusCode, output)
}
