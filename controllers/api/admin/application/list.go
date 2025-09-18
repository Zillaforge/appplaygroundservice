package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/controllers/api"
	adminCom "AppPlaygroundService/controllers/api/admin/common"
	"AppPlaygroundService/utility"
	"net/http"

	tkErr "github.com/Zillaforge/toolkits/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	cCnt "github.com/Zillaforge/appplaygroundserviceclient/constants"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

type (
	// ListApplicationsInput ...
	ListApplicationsInput struct {
		adminCom.Pagination
		Where []string `json:"where" form:"where"`
		_     struct{}
	}
	// ListApplicationsOutput ...
	ListApplicationsOutput struct {
		Applications []adminCom.ApplicationInterface `json:"applications"`
		Total        int                             `json:"total"`
		_            struct{}
	}
)

// ListSkeletons ...
func ListApplications(c *gin.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(c)

		input  = &ListApplicationsInput{}
		output = &ListApplicationsOutput{
			Applications: []adminCom.ApplicationInterface{},
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

	// call grpc ListApplications
	listApplicationsInput := &pb.ListWithLanguageInput{
		Limit:    int32(input.Limit),
		Offset:   int32(input.Offset),
		Where:    append(input.Where, "namespace="+c.GetString(cnt.CtxNamespace)),
		Language: c.GetString(cnt.CtxLanguage),
	}
	listApplicationsOutput, listApplicationsErr := aps.ListApplications(listApplicationsInput, c)
	if listApplicationsErr != nil {
		if e, ok := tkErr.IsError(listApplicationsErr); ok {
			switch e.Code() {
			case cCnt.GRPCWhereBindingErrCode:
				statusCode = http.StatusBadRequest
				err = tkErr.New(cnt.ControllerWhereQueryInvalidErr)
				utility.ResponseWithType(c, statusCode, err)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Controller, "aps.ListApplications(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listApplicationsInput),
		).Error(listApplicationsErr.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.AdminAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	for _, data := range listApplicationsOutput.Data {
		a := &adminCom.Application{}
		a.ExtractByProto(c, data)
		output.Applications = append(output.Applications, a)
	}
	output.Total = int(listApplicationsOutput.Count)

	utility.ResponseWithType(c, statusCode, output)
}
