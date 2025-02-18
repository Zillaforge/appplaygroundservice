package module

import (
	"AppPlaygroundService/modules/quiz"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"encoding/json"
	"time"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.ModuleCRUDControllerServer
type Method struct {
	// Embed UnsafeModuleCRUDControllerServer to have mustEmbedUnimplementedModuleCRUDControllerServer()
	pb.UnsafeModuleCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ModuleCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.Module) (output *pb.ModuleDetail) {
	questions := []byte{}

	output = &pb.ModuleDetail{
		Module: &pb.ModuleInfo{
			ID:               input.ID,
			Name:             input.Name,
			Description:      input.Description,
			ModuleCategoryID: input.ModuleCategoryID,
			Questions:        questions,
			Location:         input.Location,
			State:            input.State,
			Public:           input.Public,
			CreatorID:        input.CreatorID,
			CreatedAt:        input.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:        input.UpdatedAt.UTC().Format(time.RFC3339),
		},
		ModuleCategory: &pb.ModuleCategoryInfo{
			ID:          input.ModuleCategory.ID,
			Name:        input.ModuleCategory.Name,
			Description: input.ModuleCategory.Description,
			CreatorID:   input.ModuleCategory.CreatorID,
			CreatedAt:   input.ModuleCategory.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   input.ModuleCategory.UpdatedAt.UTC().Format(time.RFC3339),
		},
	}

	return
}

func (m Method) getQuestions(id string, language string) (questions []byte, err error) {
	questionInput := &quiz.ParseQuestionFileInput{
		QuestionsFolder: utility.FormatModulePath(id),
		LanguageCode:    language,
	}
	questionOutput, err := quiz.ParseQuestionFile(questionInput)
	if err != nil {
		return []byte{}, err
	}
	// Marshal the Questions object into JSON bytes
	questions, err = json.Marshal(questionOutput.ParsedQuestions)
	if err != nil {
		return []byte{}, err
	}

	return questions, nil
}
