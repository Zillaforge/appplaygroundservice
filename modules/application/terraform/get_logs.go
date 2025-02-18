package terraform

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/utility"
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func (p *TerraformProvider) GetLogs(ctx context.Context, input common.GetLogsInput) (output common.GetLogsOutput, err error) {
	output.Logs = ""

	logFilePath := filepath.Join(
		utility.FormatApplicationPath(input.ProjectID, input.ApplicationID),
		tfErrLogFile,
	)

	// check file
	fileInfo, _err := os.Stat(logFilePath)
	if _err != nil {
		zap.L().With(
			zap.String(cnt.Module, "os.Stat(...)"),
			zap.String("filePath", logFilePath),
		).Warn(_err.Error())
		return
	}
	logsByte := make([]byte, fileInfo.Size())

	// open and read log files into logsByte
	fh, err := os.Open(logFilePath)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "os.Open(...)"),
			zap.String("filePath", logFilePath),
		).Error(err.Error())
		return
	}

	defer fh.Close()

	_, err = fh.Read(logsByte[:])
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "fh.Read(...)"),
			zap.String("filePath", logFilePath),
		).Error(err.Error())
		return
	}

	output.Logs = string(logsByte)

	return
}
