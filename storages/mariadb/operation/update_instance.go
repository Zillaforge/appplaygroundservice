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

// UpdateInstance ...
func (o *Operation) UpdateInstance(ctx context.Context, input *common.UpdateInstanceInput) (output *common.UpdateInstanceOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"error":  &err,
	})

	whereCondition := &common.UpdateInstanceInput{
		ID: input.ID,
	}
	output = &common.UpdateInstanceOutput{}
	if updateErr := o.conn.WithContext(ctx).Model(&tables.Instance{}).Where(queryConversion(*whereCondition)).Updates(queryConversion(*input.UpdateData)).Preload("Application").First(&output.Instance).Error; err != nil {
		if errors.Is(updateErr, gorm.ErrRecordNotFound) {
			// Instance not found
			zap.L().With(
				zap.String(cnt.Storage, "record not found"),
				zap.String(cnt.RequestID, requestID),
				zap.String("name", GetTableName(&tables.Instance{})),
			).Error(updateErr.Error())
			err = tkErr.New(cnt.StorageInstanceNotFoundErr).WithInner(updateErr)
			return
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).Updates(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Instance{})),
			zap.Any("value", input),
		).Error(updateErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(updateErr)
		return
	}
	return
}
