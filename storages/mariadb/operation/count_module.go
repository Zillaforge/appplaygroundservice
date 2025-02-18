package operation

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// CountModule ...
func (o *Operation) CountModule(ctx context.Context, input *common.CountModuleInput) (output *common.CountModuleOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  input,
		"output": output,
		"error":  &err,
	})

	output = &common.CountModuleOutput{}
	if countModuleErr := o.conn.WithContext(ctx).Model(&tables.Module{}).Where("module.id = ?", input.ID).Count(&output.Count).Error; countModuleErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.Model(...).Where(...).Count(...).Error"),
			zap.Any(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Module{})),
		).Error(countModuleErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(countModuleErr)
		return
	}
	return
}
