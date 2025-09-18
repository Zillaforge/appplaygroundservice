package instance

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
	// ListInstancesInput ...
	ListInstancesInput struct {
		adminCom.Pagination
		Where []string `json:"where" form:"where"`
		_     struct{}
	}
	// ListInstancesOutput ...
	ListInstancesOutput struct {
		Instances []adminCom.InstanceInterface `json:"instances"`
		Total     int                          `json:"total"`
		_         struct{}
	}
)

// ListInstances ...
func ListInstances(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input  = &ListInstancesInput{}
		output = &ListInstancesOutput{
			Instances: []adminCom.InstanceInterface{},
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

	listInstancesInput := &pb.ListInput{
		Limit:  int32(input.Limit),
		Offset: int32(input.Offset),
		Where:  input.Where,
	}
	listInstancesOutput, listInstancesErr := aps.ListInstances(listInstancesInput, c)
	if listInstancesErr != nil {
		if e, ok := tkErr.IsError(listInstancesErr); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.ControllerWhereQueryInvalidErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.ListInstances(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listInstancesInput),
		).Error(listInstancesErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	for _, data := range listInstancesOutput.Data {
		a := &adminCom.Instance{}
		a.ExtractByProto(data)
		output.Instances = append(output.Instances, a)
	}

	output.Total = int(listInstancesOutput.Count)
	utility.ResponseWithType(c, statusCode, output)
}
