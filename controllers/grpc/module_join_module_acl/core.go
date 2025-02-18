package modulejoinmoduleacl

import (
	"AppPlaygroundService/modules/quiz"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"encoding/json"
	"time"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

// Method is implement all methods as pb.ModuleJoinModuleAclCRUDControllerServer
type Method struct {
	// Embed UnsafeModuleJoinModuleAclCRUDControllerServer to have mustEmbedUnimplementedModuleJoinModuleAclCRUDControllerServer()
	pb.UnsafeModuleJoinModuleAclCRUDControllerServer
}

// Verify interface compliance at compile time where appropriate
var _ pb.ModuleJoinModuleAclCRUDControllerServer = (*Method)(nil)

func (m Method) storage2pb(input *tables.ModuleJoinModuleAcl) (output *pb.ModuleJoinModuleAclInfo) {
	questions := []byte{}

	return &pb.ModuleJoinModuleAclInfo{
		ModuleCategoryID:          input.ModuleCategoryID,
		ModuleCategoryName:        input.ModuleCategoryName,
		ModuleCategoryDescription: input.ModuleCategoryDescription,
		ModuleCategoryCreatorID:   input.ModuleCategoryCreatorID,
		ModuleID:                  input.ModuleID,
		ModuleName:                input.ModuleName,
		ModuleDescription:         input.ModuleDescription,
		Questions:                 questions,
		Location:                  input.Location,
		State:                     input.State,
		Public:                    input.Public,
		ModuleCreatorID:           input.ModuleCreatorID,
		ModuleAclID:               input.ModuleAclID,
		AllowProjectID:            input.AllowProjectID,
		ModuleCategoryCreatedAt:   input.ModuleCategoryCreatedAt.UTC().Format(time.RFC3339),
		ModuleCategoryUpdatedAt:   input.ModuleCategoryUpdatedAt.UTC().Format(time.RFC3339),
		ModuleCreatedAt:           input.ModuleCreatedAt.UTC().Format(time.RFC3339),
		ModuleUpdatedAt:           input.ModuleUpdatedAt.UTC().Format(time.RFC3339),
	}
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
