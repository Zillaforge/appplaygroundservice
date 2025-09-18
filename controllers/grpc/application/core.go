package application

import (
	cnt "AppPlaygroundService/constants"
	appCom "AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/modules/fsmhandler"
	fsmCom "AppPlaygroundService/modules/fsmhandler/common/application"
	"AppPlaygroundService/modules/quiz"
	"AppPlaygroundService/storages/tables"
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.ApplicationCRUDControllerServer
type Method struct {
	// Embed UnsafeApplicationCRUDControllerServer to have mustEmbedUnimplementedApplicationCRUDControllerServer()
	pb.UnsafeApplicationCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ApplicationCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.Application) (output *pb.ApplicationDetail) {
	output = &pb.ApplicationDetail{
		Application: &pb.ApplicationInfo{
			ID:          input.ID,
			Name:        input.Name,
			Description: input.Description,
			ModuleID:    input.ModuleID,
			State:       input.State,
			Answers:     input.Answers,
			Namespace:   input.Namespace,
			Shiftable:   input.Shiftable,
			ProjectID:   input.ProjectID,
			CreatorID:   input.CreatorID,
			UpdaterID:   input.UpdaterID,
			Extra:       input.Extra,
			CreatedAt:   input.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   input.UpdatedAt.UTC().Format(time.RFC3339),
		},
		ModuleCategory: &pb.ModuleCategoryInfo{
			ID:          input.Module.ModuleCategory.ID,
			Name:        input.Module.ModuleCategory.Name,
			Description: input.Module.ModuleCategory.Description,
			CreatorID:   input.Module.ModuleCategory.CreatorID,
			CreatedAt:   input.Module.ModuleCategory.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   input.Module.ModuleCategory.UpdatedAt.UTC().Format(time.RFC3339),
		},
		Module: &pb.ModuleInfo{
			ID:               input.Module.ID,
			Name:             input.Module.Name,
			Description:      input.Module.Description,
			ModuleCategoryID: input.Module.ModuleCategoryID,
			Location:         input.Module.Location,
			State:            input.Module.State,
			CreatorID:        input.Module.CreatorID,
			CreatedAt:        input.Module.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:        input.Module.UpdatedAt.UTC().Format(time.RFC3339),
		},
	}

	return
}

func (m Method) getAnswers(ansStr string, moduleID string, language string) (answers []byte, err error) {
	parseInput := &quiz.ParseAnswerInput{
		RawStr:       ansStr,
		ModuleID:     moduleID,
		LanguageCode: language,
	}
	parseOutput, err := quiz.ParseAnswer(parseInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "question.ParseAnswer()"),
			zap.Any("input", parseInput),
		).Error(err.Error())
		return
	}

	answers, err = json.Marshal(parseOutput.ParsedAnswers)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "json.Marshal()"),
			zap.Any("v", parseOutput.ParsedAnswers),
		).Error(err.Error())
		return
	}

	return
}

func approveApplication(c context.Context, app *tables.Application) (err error) {
	config := make(map[string]interface{})
	validAnswers := quiz.Answers{}
	err = json.Unmarshal(app.Answers, &validAnswers)
	if err != nil {
		return
	}
	for _, ans := range validAnswers.Answers {
		config[ans.Variable] = ans.Value
	}

	deployInput := &appCom.DeployInput{
		ModuleID:      app.ModuleID,
		ApplicationID: app.ID,
		ProjectID:     app.ProjectID,
		UserID:        app.CreatorID,
		AppName:       app.Name,
		Namespace:     app.Namespace,
		Config:        config,
	}
	go fsmhandler.Application.Event(c, app.ID, fsmCom.ApproveEvent, deployInput)
	return
}
