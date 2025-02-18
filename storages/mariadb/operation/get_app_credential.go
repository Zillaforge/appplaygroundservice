package operation

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// GetAppCredential ...
func (o *Operation) GetAppCredential(ctx context.Context, input *common.GetAppCredentialInput) (output *common.GetAppCredentialOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	output = &common.GetAppCredentialOutput{}
	if getErr := o.conn.WithContext(ctx).Model(&tables.AppCredential{}).Where(queryConversion(*input)).First(&output.AppCredential).Error; getErr != nil {
		if errors.Is(getErr, gorm.ErrRecordNotFound) {
			zap.L().With(
				zap.String(cnt.Storage, "record not found"),
				zap.String(cnt.RequestID, requestID),
				zap.String("name", GetTableName(&tables.AppCredential{})),
				zap.Any("value", input),
			).Warn(getErr.Error())
			err = tkErr.New(cnt.StorageAppCredentialNotFoundErr).WithInner(getErr)
			return
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).First(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.AppCredential{})),
			zap.Any("value", input),
		).Error(getErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(getErr)
		return
	}
	return output, nil
}
