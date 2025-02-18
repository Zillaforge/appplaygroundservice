package application

import (
	auth "AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	cnt "AppPlaygroundService/constants"
	apsApplication "AppPlaygroundService/modules/fsmhandler/common/application"
	"AppPlaygroundService/modules/lbmevents"
	"AppPlaygroundService/modules/quiz"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/littlebell"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (m *Method) CreateApplication(ctx context.Context, input *pb.CreateApplicationInput) (output *pb.ApplicationDetail, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	// check the module exist
	getModuleInput := &storCom.GetModuleInput{
		ID: input.Application.ModuleID,
	}
	_, getModuleErr := storages.Use().GetModule(ctx, getModuleInput)
	if getModuleErr != nil {
		if e, ok := tkErr.IsError(getModuleErr); ok {
			switch e.Code() {
			case cnt.StorageModuleNotFoundErrCode:
				err = tkErr.New(cCnt.GRPCModuleNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetModule()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getModuleInput),
		).Error(getModuleErr.Error())
		err = getModuleErr
		return
	}
	// check the project exist
	getProjectInput := &storCom.GetProjectInput{
		ID: input.Application.ProjectID,
	}
	if _, getProjectErr := storages.Use().GetProject(ctx, getProjectInput); getProjectErr != nil {
		if e, ok := tkErr.IsError(getProjectErr); ok {
			switch e.Code() {
			case cnt.StorageProjectNotFoundErrCode:
				err = tkErr.New(cCnt.GRPCProjectNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetProject()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getProjectInput),
		).Error(getProjectErr.Error())
		err = getProjectErr
		return
	}

	answers := map[string]interface{}{}
	err = json.Unmarshal(input.Application.Answers, &answers)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "json.Marshal()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("data", input.Application.Answers),
		).Error(err.Error())
		return
	}

	// verify answer
	validateInput := &quiz.ValidateAnswerInput{
		QuestionsFolder: utility.FormatModulePath(input.Application.ModuleID),
		Answer:          answers,
		UserID:          input.Application.CreatorID,
		ProjectID:       input.Application.ProjectID,
	}
	validateOutput, err := quiz.ValidateAnswer(ctx, validateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.QuizOpskResourceNotFoundErrCode, cnt.QuizAnswerCannotBeEmptyErrCode, cnt.QuizAnswerTypeErrCode,
				cnt.QuizQuestionNotExistErrCode, cnt.QuizAnswerNotInOptionsErrCode, cnt.QuizUnknownQuestionTypeErrCode,
				cnt.QuizRequiredQuestionNotAnsweredErrCode, cnt.QuizOpskResourceDuplicateErrCode,
				cnt.QuizOpskResourceIsNotAvailableErrCode:
				err = tkErr.New(cCnt.GRPCQuizModuleErr, e.Message())
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "question.Validate(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", validateInput),
		).Error(err.Error())
		return
	}

	validateAnswers, err := json.Marshal(validateOutput.ValidAnswers)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "json.Marshal()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", validateOutput.ValidAnswers),
		).Error(err.Error())
		return
	}

	createInput := &storCom.CreateApplicationInput{
		Application: tables.Application{
			Name:        input.Application.Name,
			Description: input.Application.Description,
			ModuleID:    input.Application.ModuleID,
			State:       apsApplication.ReviewState,
			Answers:     validateAnswers,
			Namespace:   input.Application.Namespace,
			Shiftable:   false,
			ProjectID:   input.Application.ProjectID,
			CreatorID:   input.Application.CreatorID,
			UpdaterID:   input.Application.UpdaterID,
		},
	}

	createOutput, err := storages.Use().CreateApplication(ctx, createInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageApplicationExistErrCode:
				err = tkErr.New(cCnt.GRPCApplicationExistErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().CreateApplication()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(err.Error())
		return
	}

	getMemInput := &authCom.GetMembershipInput{
		ProjectId: input.Application.ProjectID,
		UserId:    input.Application.CreatorID,
	}
	getMemOutput, err := auth.Use().GetMembership(ctx, getMemInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "auth.Use().GetMembership()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getMemInput),
		).Error(err.Error())
		return
	}

	isTenantAdmin := false
	switch getMemOutput.TenantRole {
	case cnt.TenantAdmin.String(), cnt.TenantOwner.String():
		isTenantAdmin = true
	}

	// When the following conditions are met, the review is automatically passed
	if !mviper.GetBool("app_playground_service.scopes.enable_resource_review") || isTenantAdmin ||
		(!input.UnderReview && !validateOutput.HasGPU) {
		eventCtx := context.Background()
		eventCtx = context.WithValue(eventCtx, cnt.RequestID, requestID)
		err = approveApplication(eventCtx, &createOutput.Application)
		if err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "approveApplication()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("app", &createOutput.Application),
			).Error(err.Error())
			return
		}
	} else {
		// 取得 Project Name
		getProjectInput := &authCom.GetProjectInput{
			ID: createOutput.Application.ProjectID,
		}
		getProjectOutput, getProjectErr := auth.Use().GetProject(ctx, getProjectInput)
		if getProjectErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "auth.Use().GetProject()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getProjectInput),
			).Error(getProjectErr.Error())
			return nil, tkErr.New(cCnt.GRPCInternalServerErr)
		}

		// 取得 User Name
		getUserInput := &authCom.GetUserInput{
			ID: createOutput.Application.CreatorID,
		}
		getUserOutput, getUserErr := auth.Use().GetUser(ctx, getUserInput)
		if getUserErr != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "auth.Use().GetUser()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getUserInput),
			).Error(getUserErr.Error())
			return nil, tkErr.New(cCnt.GRPCInternalServerErr)
		}

		// 事件觸發發送給 Review Application LBM
		eventBody := &lbmevents.ReviewApplicationEvent{}
		eventBody.With(lbmevents.ReviewApplication{
			ID:                   createOutput.Application.ID,
			Name:                 createOutput.Application.Name,
			Description:          createOutput.Application.Description,
			ModuleID:             createOutput.Application.ModuleID,
			State:                createOutput.Application.State,
			Answer:               createOutput.Application.Answers,
			Namespace:            createOutput.Application.Namespace,
			Shiftable:            createOutput.Application.Shiftable,
			ProjectID:            createOutput.Application.ProjectID,
			ProjectName:          getProjectOutput.DisplayName,
			UserID:               createOutput.Application.CreatorID,
			UserName:             getUserOutput.DisplayName,
			UpdaterID:            createOutput.Application.UpdaterID,
			CreatedAt:            createOutput.Application.CreatedAt.Format(time.RFC3339),
			UpdatedAt:            createOutput.Application.UpdatedAt.Format(time.RFC3339),
			AvailabilityDistrict: mviper.GetString("app_playground_service.scopes.availability_district"),
		})
		littlebell.Publish(ctx, &littlebell.LittleBellPublishInput{
			Target: createOutput.Application.ProjectID,
			Event:  eventBody,
		})
	}

	// for get module info and module category info
	getApplicationInput := &storCom.GetApplicationInput{ID: createOutput.Application.ID}
	getApplicationOut, getApplicationErr := storages.Use().GetApplication(ctx, getApplicationInput)
	if getApplicationErr != nil {
		if e, ok := tkErr.IsError(getApplicationErr); ok {
			switch e.Code() {
			case cnt.StorageApplicationNotFoundErrCode:
				err = tkErr.New(cCnt.GRPCApplicationNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetApplication()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", createInput),
		).Error(getApplicationErr.Error())
		err = getApplicationErr
		return
	}

	if input.Application.Language != "" {
		getApplicationOut.Application.Answers, err = m.getAnswers(
			string(getApplicationOut.Application.Answers),
			getApplicationOut.Application.ModuleID,
			input.Application.Language)
		if err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "m.getAnswers()"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("ansStr", string(getApplicationOut.Application.Answers)),
				zap.String("language", input.Application.Language),
			).Error(err.Error())
			return
		}
	}
	output = m.storage2pb(&getApplicationOut.Application)

	return
}
