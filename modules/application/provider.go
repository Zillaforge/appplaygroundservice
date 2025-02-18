package application

import (
	"AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/modules/application/terraform"
	"fmt"

	"pegasus-cloud.com/aes/toolkits/mviper"
)

var provider common.Provider

func New(kind string) {
	switch kind {
	case terraform.TerraformKind:
		provider = terraform.New(&terraform.Config{
			BinPath:          mviper.GetString("app_playground_service.application.terraform.bin_path"),
			SummaryLogLength: mviper.GetInt("app_playground_service.application.summary_log_length"),
		})
	default:
		panic(fmt.Sprintf("not support application provider [%s]", kind))
	}
}

func Use() common.Provider {
	return provider
}
