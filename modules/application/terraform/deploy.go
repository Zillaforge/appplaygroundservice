package terraform

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/modules/opstkidentity"
	"AppPlaygroundService/utility"
	"bytes"
	"context"
	"os"
	"path"
	"text/template"

	"go.uber.org/zap"
	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

func (p *TerraformProvider) Deploy(ctx context.Context, input common.DeployInput) (output common.DeployOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	_, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"error":  &err,
	})

	appDir := utility.FormatApplicationPath(input.ProjectID, input.ApplicationID)
	if err = os.MkdirAll(appDir, 0770); err != nil {
		return
	}

	// Parse configs to terraform.tfvars
	moduleDir := utility.FormatModulePath(input.ModuleID)
	baseTemplate, err := template.ParseFiles(path.Join(moduleDir, tfVarsTmplFile))
	if err != nil {
		return
	}

	getAppCredInput := &pb.GetAppCredentialInput{
		UserID:    input.UserID,
		ProjectID: input.ProjectID,
		Namespace: input.Namespace,
	}
	appCred, err := aps.GetAppCredential(getAppCredInput, ctx)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "aps.GetAppCredential"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getAppCredInput),
		).Error(err.Error())
		return
	}

	opstkPID, err := opstkidentity.Use().GetOpstkPID(ctx, input.ProjectID)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "opstkidentity.Use().GetOpstkPID"),
			zap.String(cnt.RequestID, requestID),
			zap.String("projectID", input.ProjectID),
		).Error(err.Error())
		return
	}

	input.Config[TFVarsAnchorAppCredentialID] = appCred.ID
	input.Config[TFVarsAnchorAppCredentialSecret] = appCred.Secret
	input.Config[TFVarsAnchorProjectID] = opstkPID
	input.Config[TFVarsAnchorAppName] = input.AppName

	var buf bytes.Buffer
	err = baseTemplate.Execute(&buf, input.Config)
	if err != nil {
		return
	}

	name := path.Join(appDir, tfVarsFile)
	if err = os.WriteFile(name, buf.Bytes(), 0644); err != nil {
		zap.L().With(
			zap.String(cnt.Module, "os.WriteFile"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", name),
		).Error(err.Error())
		return
	}

	// Initialize Terraform
	tf, err := initTerraform(moduleDir, p.BinPath)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "initTerraform"),
			zap.String(cnt.RequestID, requestID),
			zap.String("dir", moduleDir),
			zap.String("path", p.BinPath),
		).Error(err.Error())
		return
	}

	// Deploy Terraform
	deployInput := &terraformDeployInput{
		stateFile:  path.Join(appDir, tfStateFile),
		varsFile:   path.Join(appDir, tfVarsFile),
		errLogFile: path.Join(appDir, tfErrLogFile),
	}
	if err = deploy(tf, deployInput); err != nil {
		return
	}

	tfOutput, err := getInstanceInfo(tf, deployInput.stateFile)
	if err != nil {
		return
	}

	return common.DeployOutput{
		Data: tfOutput,
	}, nil
}
