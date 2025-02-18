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

/*
UpdateMetering 更新資料庫的 Metering Record，可更新的欄位包含 EndedAt 與 LastPublishedAt

errors:
- 15000000(internal server error)
- 15000015 (application not found)
*/
func (o *Operation) UpdateMetering(ctx context.Context, input *common.UpdateMeteringInput) (output *common.UpdateMeteringOutput, err error) {
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

	output = &common.UpdateMeteringOutput{}
	if updateErr := o.conn.WithContext(ctx).Model(&tables.Metering{}).Where("application_id = ?", input.ApplicationID).Updates(queryConversion(*input.UpdateData)).First(&output.Metering).Error; updateErr != nil {
		if errors.Is(updateErr, gorm.ErrRecordNotFound) {
			// Application not found
			zap.L().With(
				zap.String(cnt.Storage, "record not found"),
				zap.String(cnt.RequestID, requestID),
				zap.String("application id", input.ApplicationID),
			).Error(updateErr.Error())
			return nil, tkErr.New(cnt.StorageApplicationNotFoundErr)
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).Updates(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Metering{})),
			zap.Any("value", input),
		).Error(updateErr.Error())
		return nil, tkErr.New(cnt.StorageInternalServerErr).WithInner(updateErr)
	}
	return output, nil
}
