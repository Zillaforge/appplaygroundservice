package constants

import tkErr "github.com/Zillaforge/toolkits/errors"

const (
	// 1600xxxx: Module
	ModuleInternalServerErrCode = 16000000
	ModuleInternalServerErrMsg  = "internal server error"

	// 1601xxxx: OpskResource
	OpskResourceRecordNotFoundErrCode = 16010000
	OpskResourceRecordNotFoundErrMsg  = "record not found"

	// 1602xxxx: Quiz
	QuizOpskResourceNotFoundErrCode        = 16020000
	QuizOpskResourceNotFoundErrMsg         = "%s (%s) not found"
	QuizAnswerCannotBeEmptyErrCode         = 16020001
	QuizAnswerCannotBeEmptyErrMsg          = "answer cannot be empty: %s"
	QuizAnswerTypeErrCode                  = 16020002
	QuizAnswerTypeErrMsg                   = "answer type error: %s"
	QuizQuestionNotExistErrCode            = 16020003
	QuizQuestionNotExistErrMsg             = "question not exist: %s"
	QuizAnswerNotInOptionsErrCode          = 16020004
	QuizAnswerNotInOptionsErrMsg           = "answer not in options: %s"
	QuizUnknownQuestionTypeErrCode         = 16020005
	QuizUnknownQuestionTypeErrMsg          = "unknown question type: %s"
	QuizRequiredQuestionNotAnsweredErrCode = 16020006
	QuizRequiredQuestionNotAnsweredErrMsg  = "required question not answered: %s"
	QuizOpskResourceDuplicateErrCode       = 16020007
	QuizOpskResourceDuplicateErrMsg        = "duplicate %s (%s)"
	QuizOpskResourceIsNotAvailableErrCode  = 16020008
	QuizOpskResourceIsNotAvailableErrMsg   = "%s (%s) is not available"
)

var (
	// 1600xxxx: Module
	// 16000000(internal server error)
	ModuleInternalServerErr = tkErr.Error(ModuleInternalServerErrCode, ModuleInternalServerErrMsg)

	// 16010000(record not found)
	OpskResourceRecordNotFoundErr = tkErr.Error(OpskResourceRecordNotFoundErrCode, OpskResourceRecordNotFoundErrMsg)

	// 16020000(%s (%s) not found) ex. flavor (3fe78aca-3ac7-4051-a1f0-5baf3d20443f) not found
	QuizOpskResourceNotFoundErr = tkErr.Error(QuizOpskResourceNotFoundErrCode, QuizOpskResourceNotFoundErrMsg)
	// 16020001(answer cannot be empty: %s) ex. answer cannot be empty: flavor_id
	QuizAnswerCannotBeEmptyErr = tkErr.Error(QuizAnswerCannotBeEmptyErrCode, QuizAnswerCannotBeEmptyErrMsg)
	// 16020002(answer type error: %s) ex. answer type error: flavor_id
	QuizAnswerTypeErr = tkErr.Error(QuizAnswerTypeErrCode, QuizAnswerTypeErrMsg)
	// 16020003(question not exist: %s) ex. question not exist: flavor_id
	QuizQuestionNotExistErr = tkErr.Error(QuizQuestionNotExistErrCode, QuizQuestionNotExistErrMsg)
	// 16020004(answer not in options: %s) ex. answer not in options: flavor_id
	QuizAnswerNotInOptionsErr = tkErr.Error(QuizAnswerNotInOptionsErrCode, QuizAnswerNotInOptionsErrMsg)
	// 16020005(unknown question type: %s) ex. unknown question type: flavor_id
	QuizUnknownQuestionTypeErr = tkErr.Error(QuizUnknownQuestionTypeErrCode, QuizUnknownQuestionTypeErrMsg)
	// 16020006(required question not answered: %s) ex. required question not answered: flavor_id
	QuizRequiredQuestionNotAnsweredErr = tkErr.Error(QuizRequiredQuestionNotAnsweredErrCode, QuizRequiredQuestionNotAnsweredErrMsg)
	// 16020007(duplicate %s (%s)) ex. duplicate volume (3fe78aca-3ac7-4051-a1f0-5baf3d20443f)
	QuizOpskResourceDuplicateErr = tkErr.Error(QuizOpskResourceDuplicateErrCode, QuizOpskResourceDuplicateErrMsg)
	// 16020008(%s (%s) is not available) ex. volume (3fe78aca-3ac7-4051-a1f0-5baf3d20443f) is not available
	QuizOpskResourceIsNotAvailableErr = tkErr.Error(QuizOpskResourceIsNotAvailableErrCode, QuizOpskResourceIsNotAvailableErrMsg)
)
