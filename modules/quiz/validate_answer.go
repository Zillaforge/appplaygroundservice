package quiz

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/opskresource"
	opskCom "AppPlaygroundService/modules/opskresource/common"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

type ValidateAnswerInput struct {
	QuestionsFolder string
	Answer          map[string]interface{}

	UserID    string
	ProjectID string
}

type ValidateAnswerOutput struct {
	ValidAnswers Answers
	HasGPU       bool
}

// ValidateAnswer is the function to validate the answer
func ValidateAnswer(ctx context.Context, input *ValidateAnswerInput) (output *ValidateAnswerOutput, err error) {
	// Load the question.json file to get the questions
	questionPath := fmt.Sprintf("%s/questions.json", input.QuestionsFolder)
	content, err := os.ReadFile(questionPath)
	if err != nil {
		return
	}

	questions := Questions{}
	err = json.Unmarshal(content, &questions)
	if err != nil {
		return
	}

	questionMap := make(map[string]Question)
	for _, question := range questions.Questions {
		questionMap[question.Variable] = question
	}

	// Validate the answer
	validAnswer, err := validate(ctx, questionMap, input)
	if err != nil {
		return
	}

	output = &ValidateAnswerOutput{
		ValidAnswers: *validAnswer,
		HasGPU:       hasGPU(questionMap),
	}
	return
}

// validate is the function to validate the answer
func validate(ctx context.Context, question map[string]Question, input *ValidateAnswerInput) (validAnswer *Answers, err error) {
	answer := input.Answer

	questionMap := make(map[string]Question)
	for key, value := range question {
		questionMap[key] = value
	}

	for key, value := range answer {
		if _, ok := questionMap[key]; !ok {
			err = tkErr.New(cnt.QuizQuestionNotExistErr, key)
			return
		}

		// check type format
		var ok bool
		switch questionMap[key].Type {
		case TypeNameString, TypeNamePassword, TypeNameInt, TypeNameVPSNetwork, TypeNameVPSFlavor,
			TypeNameVPSGPUFlavor, TypeNameVPSvGPUFlavor, TypeNameVPSKeypair, TypeNameVPSVolume,
			TypeNameVPSBootVolume, TypeNamePort:
			var v string
			if v, ok = value.(string); !ok {
				err = tkErr.New(cnt.QuizAnswerTypeErr, key)
				return
			}
			if v == "" && questionMap[key].Required {
				err = tkErr.New(cnt.QuizAnswerCannotBeEmptyErr, key)
				return
			}
		case TypeNameArray, TypeNameVPSSecurityGroups:
			var v []string
			if v, ok = value.([]string); !ok {
				if interfaceSlice, ok := value.([]interface{}); ok {
					for _, item := range interfaceSlice {
						if str, ok := item.(string); ok {
							v = append(v, str)
						} else {
							err = tkErr.New(cnt.QuizAnswerTypeErr, key)
							return
						}
					}
				} else {
					err = tkErr.New(cnt.QuizAnswerTypeErr, key)
					return
				}
			}
			if len(v) == 0 {
				err = tkErr.New(cnt.QuizAnswerCannotBeEmptyErr, key)
				return
			}
		case TypeNameBoolean:
			if _, ok = value.(bool); !ok {
				err = tkErr.New(cnt.QuizAnswerTypeErr, key)
				return
			}
		case TypeNameEnum:
			var v string
			if v, ok = value.(string); !ok {
				err = tkErr.New(cnt.QuizAnswerTypeErr, key)
				return
			}
			if v == "" && questionMap[key].Required {
				err = tkErr.New(cnt.QuizAnswerCannotBeEmptyErr, key)
				return
			}
			// Check enum type must be in the options
			if questionMap[key].Options != nil {
				var found bool
				for _, v := range *questionMap[key].Options {
					if v == value {
						found = true
						break
					}
				}
				if !found {
					err = tkErr.New(cnt.QuizAnswerNotInOptionsErr, key)
					return
				}
			}
		case TypeNameSSHPort: // allow empty to disable ssh
			if _, ok = value.(string); !ok {
				err = tkErr.New(cnt.QuizAnswerTypeErr, key)
				return
			}
		default:
			err = tkErr.New(cnt.QuizUnknownQuestionTypeErr, key)
			return
		}
		delete(questionMap, key)
	}

	// Not answered
	for key, q := range questionMap {
		if q.Required {
			err = tkErr.New(cnt.QuizRequiredQuestionNotAnsweredErr, key)
			return
		} else {
			// Fill the default value
			answer[key] = q.Default
		}
	}

	// check VPS resource and fill the answer value into the Answer struct
	validAnswer = &Answers{}
	duplicateVolume := map[string]bool{}
	duplicateSG := map[string]bool{}
	for key, ans := range answer {
		// check VPS resource
		var displayName, rawValue, value interface{}
		// keep raw value
		rawValue = ans
		var ok bool
		switch question[key].Type {
		case TypeNameVPSFlavor, TypeNameVPSGPUFlavor, TypeNameVPSvGPUFlavor:
			v := ans.(string)
			opskGetFlavorInput := &opskCom.GetFlavorInput{ID: v}
			opskGetFlavorOutput, _err := opskresource.Use().GetFlavor(ctx, opskGetFlavorInput)
			if _err != nil {
				if e, ok := tkErr.IsError(_err); ok {
					switch e.Code() {
					case cnt.OpskResourceRecordNotFoundErrCode:
						err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "flavor", v)
						return
					}
				}
				err = _err
				return
			}
			value = ans
			displayName = opskGetFlavorOutput.Name
		case TypeNameVPSNetwork:
			v := ans.(string)
			opskGetNetworkInput := &opskCom.GetNetworkInput{ID: v}
			opskGetNetworkOutput, _err := opskresource.Use().GetNetwork(ctx, opskGetNetworkInput)
			if _err != nil {
				if e, ok := tkErr.IsError(_err); ok {
					switch e.Code() {
					case cnt.OpskResourceRecordNotFoundErrCode:
						err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "network", v)
						return
					}
				}
				err = _err
				return
			}

			// check resource permission
			if opskGetNetworkOutput.ProjectID != input.ProjectID {
				err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "network", v)
				return
			}

			value = ans
			displayName = opskGetNetworkOutput.Name
		case TypeNameVPSKeypair:
			v := ans.(string)
			opskGetKeypairInput := &opskCom.GetKeypairInput{ID: v}
			opskGetKeypairOutput, _err := opskresource.Use().GetKeypair(ctx, opskGetKeypairInput)
			if _err != nil {
				if e, ok := tkErr.IsError(_err); ok {
					switch e.Code() {
					case cnt.OpskResourceRecordNotFoundErrCode:
						err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "keypair", v)
						return
					}
				}
				err = _err
				return
			}

			// check resource permission
			if opskGetKeypairOutput.UserID != input.UserID {
				err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "keypair", v)
				return
			}

			value = opskGetKeypairOutput.Name
			displayName = opskGetKeypairOutput.Name
		case TypeNameVPSVolume:
			v := ans.(string)
			// check duplicate volume
			if _, ok = duplicateVolume[v]; !ok {
				duplicateVolume[v] = true
			} else {
				err = tkErr.New(cnt.QuizOpskResourceDuplicateErr, "volume", v)
				return
			}

			opskGetVolumeInput := &opskCom.GetVolumeInput{ID: v}
			opskGetVolumeOutput, _err := opskresource.Use().GetVolume(ctx, opskGetVolumeInput)
			if _err != nil {
				if e, ok := tkErr.IsError(_err); ok {
					switch e.Code() {
					case cnt.OpskResourceRecordNotFoundErrCode:
						err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "volume", v)
						return
					}
				}
				err = _err
				return
			}

			// check resource permission
			if opskGetVolumeOutput.ProjectID != input.ProjectID {
				err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "volume", v)
				return
			}

			// check if volume available
			if opskGetVolumeOutput.Status != "available" {
				err = tkErr.New(cnt.QuizOpskResourceIsNotAvailableErr, "volume", v)
				return
			}

			value = ans
			displayName = opskGetVolumeOutput.Name
		case TypeNameVPSSecurityGroups:
			var v, names []string
			if v, ok = ans.([]string); !ok {
				if interfaceSlice, ok := ans.([]interface{}); ok {
					for _, item := range interfaceSlice {
						if str, ok := item.(string); ok {
							v = append(v, str)
						}
					}
				}
			}

			for _, id := range v {
				// check duplicate security group
				if _, ok = duplicateSG[id]; !ok {
					duplicateSG[id] = true
				} else {
					err = tkErr.New(cnt.QuizOpskResourceDuplicateErr, "security group", id)
					return
				}

				opskGetSecurityGroupInput := &opskCom.GetSecurityGroupInput{ID: id}
				opskGetSecurityGroupOutput, _err := opskresource.Use().GetSecurityGroup(ctx, opskGetSecurityGroupInput)
				if _err != nil {
					if e, ok := tkErr.IsError(_err); ok {
						switch e.Code() {
						case cnt.OpskResourceRecordNotFoundErrCode:
							err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "security group", id)
							return
						}
					}
					err = _err
					return
				}

				// check resource permission
				if opskGetSecurityGroupOutput.ProjectID != input.ProjectID {
					err = tkErr.New(cnt.QuizOpskResourceNotFoundErr, "security group", id)
					return
				}

				names = append(names, opskGetSecurityGroupOutput.Name)
			}
			value = ans
			displayName = names
		// password encryption
		case TypeNamePassword:
			var _err error
			v := ans.(string)
			value, _err = hashPassword(v)
			if _err != nil {
				err = fmt.Errorf("hashPassword error")
				return
			}
			displayName = value
		case TypeNameString, TypeNameInt, TypeNameEnum, TypeNameArray, TypeNameBoolean, TypeNameVPSBootVolume,
			TypeNamePort, TypeNameSSHPort:
			value = ans
			displayName = ans
		default:
			err = tkErr.New(cnt.QuizUnknownQuestionTypeErr, key)
			return
		}

		// fill answer
		validAnswer.Answers = append(validAnswer.Answers, Answer{
			Question:    question[key],
			RawValue:    rawValue,
			Value:       value,
			DisplayName: displayName,
		})
	}

	return
}

func hasGPU(question map[string]Question) bool {
	for key := range question {
		switch question[key].Type {
		case TypeNameVPSGPUFlavor, TypeNameVPSvGPUFlavor:
			return true
		}
	}
	return false
}

// hashPassword hashes the given password using bcrypt.
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
