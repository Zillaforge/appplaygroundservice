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
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

/*
GetMetering 回傳指定 ID 的 Metering Record

errors:
- 15000000(internal server error)
- 15000022(metering not found)
*/
func (o *Operation) GetMetering(ctx context.Context, input *common.GetMeteringInput) (output *common.GetMeteringOutput, err error) {
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

	output = &common.GetMeteringOutput{}
	if getErr := o.conn.WithContext(ctx).Model(&tables.Metering{}).Where("metering.application_id = ?", input.ID).First(&output.Metering).Error; getErr != nil {
		if errors.Is(getErr, gorm.ErrRecordNotFound) {
			// Metering not found
			zap.L().With(
				zap.String(cnt.Storage, "record not found"),
				zap.String(cnt.RequestID, requestID),
				zap.String("name", GetTableName(&tables.Metering{})),
			).Error(getErr.Error())
			err = tkErr.New(cnt.StorageMeteringNotFoundErr).WithInner(getErr)
			return
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).First(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Metering{})),
		).Error(getErr.Error())
		return nil, tkErr.New(cnt.StorageInternalServerErr).WithInner(getErr)
	}
	return output, nil
}
