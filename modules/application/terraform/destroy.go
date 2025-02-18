package terraform

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/utility"
	"context"
	"os"
	"path"

	"go.uber.org/zap"
)

func (p *TerraformProvider) Destroy(ctx context.Context, input common.DestroyInput) (err error) {
	appDir := utility.FormatApplicationPath(input.ProjectID, input.ApplicationID)
	moduleDir := utility.FormatModulePath(input.ModuleID)

	tf, err := initTerraform(moduleDir, p.BinPath)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "initTerraform"),
			zap.String("dir", moduleDir),
			zap.String("path", p.BinPath),
		).Error(err.Error())
		return
	}

	destroyInput := &terraformDestroyInput{
		stateFile:  path.Join(appDir, tfStateFile),
		varsFile:   path.Join(appDir, tfVarsFile),
		errLogFile: path.Join(appDir, tfErrLogFile),
	}
	if err = destroy(tf, destroyInput); err != nil {
		zap.L().With(
			zap.String(cnt.Module, "destroy"),
			zap.Any("input", destroyInput),
		).Error(err.Error())
		return
	}

	// delete application directory
	os.RemoveAll(appDir)

	return
}
