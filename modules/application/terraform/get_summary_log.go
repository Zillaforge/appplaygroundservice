package terraform

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/utility"
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

func (p *TerraformProvider) GetSummaryLog(ctx context.Context, input common.GetSummaryLogInput) (output common.GetSummaryLogOutput, err error) {
	output.Log = ""

	logFilePath := filepath.Join(
		utility.FormatApplicationPath(input.ProjectID, input.ApplicationID),
		tfErrLogFile,
	)

	// open and read log files into logByte
	file, err := os.Open(logFilePath)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "os.Open(...)"),
			zap.String("filePath", logFilePath),
		).Error(err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Error:") {
			if p.SummaryLogLength == -1 {
				output.Log = line
				return
			} else if p.SummaryLogLength > 0 && len(line) > p.SummaryLogLength {
				output.Log = line[:p.SummaryLogLength] + "..."
				return
			} else {
				err = fmt.Errorf("input length error: length = %d", p.SummaryLogLength)
				zap.L().With(
					zap.Int("length", p.SummaryLogLength),
				).Error(err.Error())
				return
			}
		}
	}

	if err = scanner.Err(); err != nil {
		zap.L().With(
			zap.String(cnt.Module, "scanner.Err()"),
		).Error(err.Error())
		return
	}

	err = fmt.Errorf("summary not found")
	zap.L().With(
		zap.String("filePath", logFilePath),
	).Error(err.Error())
	return
}
