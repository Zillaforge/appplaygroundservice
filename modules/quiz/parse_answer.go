package quiz

import (
	"AppPlaygroundService/utility"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"

	mviper "pegasus-cloud.com/aes/toolkits/mviper"
)

type ParseAnswerInput struct {
	RawStr       string
	ModuleID     string
	LanguageCode string
}

type ParseAnswerOutput struct {
	ParsedAnswers Answers
}

func ParseAnswer(input *ParseAnswerInput) (output *ParseAnswerOutput, err error) {
	baseTemplate, err := template.New("parse").Parse(input.RawStr)
	if err != nil {
		return
	}

	// Load the questions_strings_{lang_code}.yaml file
	questionsFolder := utility.FormatModulePath(input.ModuleID)
	langStringsPath := filepath.Join(questionsFolder, fmt.Sprintf(filepathTmpl, input.LanguageCode))
	langStrings, err := loadLanguageStrings(langStringsPath)
	if err != nil {
		// check if defaultLanguage file can be use
		defaultLanguage := mviper.GetString("app_playground_service.scopes.default_language")
		if defaultLanguage != input.LanguageCode {
			langStringsPath = filepath.Join(questionsFolder, fmt.Sprintf(filepathTmpl, defaultLanguage))
			langStrings, err = loadLanguageStrings(langStringsPath)
		}

		if err != nil {
			return
		}
	}

	// Parse the question.json file as template and replace the placeholders with the values from the yaml file
	var buf bytes.Buffer
	err = baseTemplate.Execute(&buf, langStrings)
	if err != nil {
		return
	}

	ans := Answers{}
	decoder := json.NewDecoder(&buf)
	err = decoder.Decode(&ans)
	if err != nil {
		return
	}

	output = &ParseAnswerOutput{
		ParsedAnswers: ans,
	}
	return
}
