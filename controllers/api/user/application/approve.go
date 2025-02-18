package application

import (
	auth "AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	cnt "AppPlaygroundService/constants"
	userCom "AppPlaygroundService/controllers/api/user/common"
	appCom "AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/modules/fsmhandler"
	fsmCom "AppPlaygroundService/modules/fsmhandler/common/application"
	"AppPlaygroundService/modules/lbmevents"
	"AppPlaygroundService/modules/quiz"
	"AppPlaygroundService/utility"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/littlebell"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type ApproveApplicationInput struct {
	ID string `json:"-"`
	_  struct{}
}

type ApproveApplicationOutput struct {
	userCom.Application
	_ struct{}
}

func ApproveApplication(c *gin.Context) {
	var (
		input        = &ApproveApplicationInput{ID: c.GetString(cnt.CtxApplicationID)}
		output       = &ApproveApplicationOutput{}
		err          error
		requestID        = utility.MustGetContextRequestID(c)
		funcName         = tkUtils.NameOfFunction().Name()
		statusCode   int = http.StatusOK
		supportRoles     = map[string]bool{
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

	// Approve can only perform by TENANT_ADMIN
	if role := c.GetString(cnt.CtxTenantRole); !supportRoles[role] &&
		c.GetString(cnt.CtxCreator) != c.GetString(cnt.CtxUserID) {
		statusCode = http.StatusUnauthorized
		err = tkErr.New(cnt.UserAPIUnauthorizedOpErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// get application
	getApplicationInput := &pb.GetWithLanguageInput{
		ID:       input.ID,
		Language: c.GetString(cnt.CtxLanguage),
	}
	getApplicationOutput, err := aps.GetApplication(getApplicationInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getApplicationInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	err = approveApplication(c, getApplicationOutput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "approveApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getApplicationOutput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	// get lastest application
	getApplicationOutput, err = aps.GetApplication(getApplicationInput, c)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "aps.GetApplication(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getApplicationInput),
		).Error(err.Error())
		statusCode = http.StatusInternalServerError
		err = tkErr.New(cnt.UserAPIInternalServerErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	output.ExtractByProto(c, getApplicationOutput)
	utility.ResponseWithType(c, statusCode, output)
}

func approveApplication(c *gin.Context, app *pb.ApplicationDetail) (err error) {
	config := make(map[string]interface{})
	validAnswers := quiz.Answers{}
	err = json.Unmarshal(app.Application.Answers, &validAnswers)
	if err != nil {
		return
	}
	for _, ans := range validAnswers.Answers {
		config[ans.Variable] = ans.Value
	}

	deployInput := &appCom.DeployInput{
		ModuleID:      app.Module.ID,
		ApplicationID: app.Application.ID,
		ProjectID:     app.Application.ProjectID,
		UserID:        app.Application.CreatorID,
		AppName:       app.Application.Name,
		Namespace:     app.Application.Namespace,
		Config:        config,
	}

	if lbmPublishErr := lbmPublish(c, app, c.GetString(cnt.CtxUserID), true); lbmPublishErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "lbmPublish()"),
			zap.Any("app", app),
			zap.String("reviewerID", c.GetString(cnt.CtxUserID)),
		).Error(lbmPublishErr.Error())
		return tkErr.New(cnt.UserAPIInternalServerErr)
	}

	go fsmhandler.Application.Event(c, app.Application.ID, fsmCom.ApproveEvent, deployInput)
	return
}

func lbmPublish(ctx context.Context, app *pb.ApplicationDetail, reviewerID string, approved bool) (err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	_, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"app":        &app,
		"reviewerID": &reviewerID,
		"approved":   &approved,
		"error":      &err,
	})

	// 取得 Project Name
	getProjectInput := &authCom.GetProjectInput{
		ID: app.Application.ProjectID,
	}
	getProjectOutput, getProjectErr := auth.Use().GetProject(ctx, getProjectInput)
	if getProjectErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "auth.Use().GetProject()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getProjectInput),
		).Error(getProjectErr.Error())
		return tkErr.New(cnt.UserAPIInternalServerErr)
	}

	// 取得 User Name
	getUserInput := &authCom.GetUserInput{
		ID: app.Application.CreatorID,
	}
	getUserOutput, getUserErr := auth.Use().GetUser(ctx, getUserInput)
	if getUserErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "auth.Use().GetUser()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getUserInput),
		).Error(getUserErr.Error())
		return tkErr.New(cnt.UserAPIInternalServerErr)
	}

	// 取得 Reviewer Name
	getReviewerInput := &authCom.GetUserInput{
		ID: reviewerID,
	}
	getReviewerOutput, getReviewerErr := auth.Use().GetUser(ctx, getReviewerInput)
	if getReviewerErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "auth.Use().GetUser()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getReviewerInput),
		).Error(getReviewerErr.Error())
		return tkErr.New(cnt.UserAPIInternalServerErr)
	}

	// 事件觸發發送 Approve Application 給 LBM
	littlebell.Publish(ctx, &littlebell.LittleBellPublishInput{
		Target: getUserOutput.ID,
		Event: func() littlebell.EventIntf {
			if approved {
				eventBody := &lbmevents.ApproveApplicationEvent{}
				eventBody.With(lbmevents.ApproveApplication{
					AvailabilityDistrict: mviper.GetString("app_playground_service.scopes.availability_district"),
					ID:                   app.Application.ID,
					Namespace:            app.Application.Namespace,
					Name:                 app.Application.Name,
					ProjectID:            app.Application.ProjectID,
					ProjectName:          getProjectOutput.DisplayName,
					ReviewerID:           reviewerID,
					ReviewerName:         getReviewerOutput.DisplayName,
					UserID:               getUserOutput.ID,
					UserName:             getUserOutput.DisplayName,
					CreatedAt:            app.Application.CreatedAt,
				})
				return eventBody
			} else {
				eventBody := &lbmevents.RejectApplicationEvent{}
				eventBody.With(lbmevents.RejectApplication{
					AvailabilityDistrict: mviper.GetString("app_playground_service.scopes.availability_district"),
					ID:                   app.Application.ID,
					Namespace:            app.Application.Namespace,
					Name:                 app.Application.Name,
					ProjectID:            app.Application.ProjectID,
					ProjectName:          getProjectOutput.DisplayName,
					ReviewerID:           reviewerID,
					ReviewerName:         getReviewerOutput.DisplayName,
					UserID:               getUserOutput.ID,
					UserName:             getUserOutput.DisplayName,
					CreatedAt:            app.Application.CreatedAt,
				})
				return eventBody
			}
		}(),
	})
	return nil
}
