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
CountMetering 回傳計量資料總筆數

errors:
- 15000000(internal server error)
*/
func (o *Operation) CountMetering(ctx context.Context, input *common.CountMeteringInput) (output *common.CountMeteringOutput, err error) {
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

	output = &common.CountMeteringOutput{}
	if countMeteringErr := o.conn.WithContext(ctx).Model(&tables.Metering{}).Count(&output.Count).Error; countMeteringErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.Model(...).Where(...).Count(...).Error"),
			zap.Any(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Metering{})),
		).Error(countMeteringErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(countMeteringErr)
		return
	}
	return
}
