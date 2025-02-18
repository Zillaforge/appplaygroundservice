package operation

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages/common"
	mariadb_com "AppPlaygroundService/storages/mariadb/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// DeleteInstance ...
func (o *Operation) DeleteInstance(ctx context.Context, input *common.DeleteInstanceInput) (output *common.DeleteInstanceOutput, err error) {
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

	var id []string
	if deleteErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.Instance{}), &input.Where).Pluck("id", &id).Delete(&tables.Instance{}).Error; deleteErr != nil {
		if sqlErr, ok := deleteErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// Instance already by reference
			case mariadb_com.ER_ROW_IS_REFERENCED_2:
				zap.L().With(
					zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
					zap.String(cnt.RequestID, requestID),
					zap.String("name", GetTableName(&tables.Instance{})),
					zap.Any("value", input),
				).Warn(deleteErr.Error())
				err = tkErr.New(cnt.StorageInstanceInUseErr).WithInner(deleteErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Pluck(...).Delete(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Instance{})),
			zap.Any("value", input),
		).Error(deleteErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(deleteErr)
		return
	}
	output = &common.DeleteInstanceOutput{
		ID: id,
	}
	return
}
