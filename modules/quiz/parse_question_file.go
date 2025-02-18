package quiz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	mviper "pegasus-cloud.com/aes/toolkits/mviper"
)

const filepathTmpl string = "questions_strings_%s.yaml"

type ParseQuestionFileInput struct {
	QuestionsFolder string
	LanguageCode    string
}

type ParseQuestionFileOutput struct {
	ParsedQuestions Questions
}

// ParseQuestionFile is the function to load the questions from the question.json file and replace the placeholders with the translation from the yaml file
func ParseQuestionFile(input *ParseQuestionFileInput) (output *ParseQuestionFileOutput, err error) {
	// Load the question.json file
	questionPath := filepath.Join(input.QuestionsFolder, "questions.json")
	baseTemplate, err := template.ParseFiles(questionPath)
	if err != nil {
		return
	}

	// Load the questions_strings_{lang_code}.yaml file
	langStringsPath := filepath.Join(input.QuestionsFolder, fmt.Sprintf(filepathTmpl, input.LanguageCode))
	langStrings, err := loadLanguageStrings(langStringsPath)
	if err != nil {
		// check if defaultLanguage file can be use
		defaultLanguage := mviper.GetString("app_playground_service.scopes.default_language")
		if defaultLanguage != input.LanguageCode {
			langStringsPath = filepath.Join(input.QuestionsFolder, fmt.Sprintf(filepathTmpl, defaultLanguage))
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

	ques := Questions{}
	decoder := json.NewDecoder(&buf)
	err = decoder.Decode(&ques)
	if err != nil {
		return
	}

	output = &ParseQuestionFileOutput{
		ParsedQuestions: ques,
	}
	return
}

func loadLanguageStrings(path string) (data map[string]interface{}, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &data)
	if err != nil {
		return
	}
	return
}
