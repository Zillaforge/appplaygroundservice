package terraform

import (
	"AppPlaygroundService/modules/application/common"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"time"

	"github.com/hashicorp/terraform-exec/tfexec"
)

const (
	TerraformKind = "terraform"

	TFVarsAnchorProjectID           = "project_id"
	TFVarsAnchorAppName             = "app_name"
	TFVarsAnchorAppCredentialID     = "application_credential_id"
	TFVarsAnchorAppCredentialSecret = "application_credential_secret"
)

type Config struct {
	BinPath          string
	SummaryLogLength int
}

type TerraformProvider struct {
	BinPath          string
	SummaryLogLength int
}

const (
	tfStateFile    = "terraform.tfstate"
	tfVarsFile     = "terraform.tfvars"
	tfVarsTmplFile = "terraform.tfvars.tmpl"
	tfErrLogFile   = "terraform_error.log"
)

type terraformDeployInput struct {
	stateFile  string
	varsFile   string
	errLogFile string
}

type terraformDestroyInput struct {
	stateFile  string
	varsFile   string
	errLogFile string
}

func New(config *Config) common.Provider {
	return &TerraformProvider{
		BinPath:          config.BinPath,
		SummaryLogLength: config.SummaryLogLength,
	}
}

func initTerraform(dir, path string) (*tfexec.Terraform, error) {
	// Initialize Terraform
	tf, err := tfexec.NewTerraform(dir, path)
	if err != nil {
		return nil, err
	}

	// Run terraform init
	err = tf.Init(context.Background())
	if err != nil {
		return nil, err
	}

	return tf, nil
}

func deploy(tf *tfexec.Terraform, input *terraformDeployInput) error {
	planSuccess, err := tf.Plan(context.Background(), tfexec.VarFile(input.varsFile), tfexec.State(input.stateFile))
	if err != nil {
		logErr := logError(input.errLogFile, fmt.Sprintf("error planning Terraform configuration: %v", err), false)
		if logErr != nil {
			return fmt.Errorf("error planning Terraform configuration: %v; additionally, failed to log error: %v", err, logErr)
		}
		return fmt.Errorf("error planning Terraform configuration: %v", err)
	}

	if !planSuccess {
		logErr := logError(input.errLogFile, "plan failed due to insufficient resources or other errors", false)
		if logErr != nil {
			return fmt.Errorf("plan failed due to insufficient resources or other errors; additionally, failed to log error: %v", logErr)
		}

		return fmt.Errorf("plan failed due to insufficient resources or other errors")
	}

	if err := tf.Apply(context.Background(), tfexec.VarFile(input.varsFile), tfexec.State(input.stateFile)); err != nil {
		logErr := logError(input.errLogFile, fmt.Sprintf("error applying Terraform configuration: %v", err), false)
		if logErr != nil {
			return fmt.Errorf("error applying Terraform configuration: %v; additionally, failed to log error: %v", err, logErr)
		}
		destroyInput := &terraformDestroyInput{
			stateFile:  input.stateFile,
			varsFile:   input.varsFile,
			errLogFile: input.errLogFile,
		}
		if destroyErr := destroy(tf, destroyInput); destroyErr != nil {
			logErr := logError(input.errLogFile, fmt.Sprintf("Destroy failed: %v", destroyErr), true)
			if logErr != nil {
				return fmt.Errorf("error applying Terraform configuration: %v; additionally, failed to log destroy error: %v", err, logErr)
			}

		}
		return fmt.Errorf("error applying Terraform configuration: %v", err)
	}

	return nil
}

func getInstanceInfo(tf *tfexec.Terraform, stateFile string) ([]common.InstanceInfo, error) {
	output, err := tf.Output(context.Background(), tfexec.State(stateFile))
	if err != nil {
		return nil, fmt.Errorf("error getting Terraform output: %v", err)
	}

	var instanceInfo []common.InstanceInfo
	if instancesMeta, ok := output["instances"]; ok {
		var instances []map[string]interface{}
		if err := json.Unmarshal(instancesMeta.Value, &instances); err != nil {
			return nil, fmt.Errorf("error unmarshaling instances: %v", err)
		}

		for _, instance := range instances {
			instance["provider"] = TerraformKind
			instanceJSON, err := json.Marshal(map[string]interface{}{
				"instance": instance,
			})
			if err != nil {
				return nil, fmt.Errorf("error marshaling instance to JSON: %v", err)
			}

			instanceInfo = append(instanceInfo, common.InstanceInfo{
				Name:        instance["name"].(string),
				ReferenceID: instance["id"].(string),
				Extra:       instanceJSON,
			})
		}
	}

	return instanceInfo, nil
}

func destroy(tf *tfexec.Terraform, input *terraformDestroyInput) error {
	if err := tf.Destroy(context.Background(), tfexec.VarFile(input.varsFile), tfexec.State(input.stateFile)); err != nil {

		logErr := logError(input.errLogFile, fmt.Sprintf("error destroying Terraform configuration: %v", err), false)

		if logErr != nil {
			return fmt.Errorf("error destroying Terraform configuration: %v; additionally, failed to log error: %v", err, logErr)
		}

		return fmt.Errorf("error destroying Terraform configuration: %v", err)
	}
	return nil
}

func logError(errLogFile, message string, append bool) error {
	var f *os.File
	var err error
	if append {
		f, err = os.OpenFile(errLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.OpenFile(errLogFile, os.O_CREATE|os.O_WRONLY, 0644)
	}

	if err != nil {
		return fmt.Errorf("error opening error log file: %v", err)
	}
	defer f.Close()

	message = fmt.Sprintf("=== %s\n%s", time.Now().Format("2006/01/02 - 15:04:05"), message)
	if _, err = f.WriteString(message + "\n"); err != nil {
		return fmt.Errorf("error writing to error log file: %v", err)
	}
	return nil
}
