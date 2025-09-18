package operation

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

/*
DeleteMetering 負責刪除指定 ID 的 Metering Record

errors:
- 15000000(internal server error)
*/
func (o *Operation) DeleteMetering(ctx context.Context, input *common.DeleteMeteringInput) (output *common.DeleteMeteringOutput, err error) {
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

	if deleteErr := o.conn.WithContext(ctx).Where("application_id = ?", input.ID).Delete(&tables.Metering{}).Error; deleteErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Delete(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Metering{})),
			zap.Any("value", input),
		).Error(deleteErr.Error())
		return nil, tkErr.New(cnt.StorageInternalServerErr).WithInner(deleteErr)
	}
	return nil, nil
}
