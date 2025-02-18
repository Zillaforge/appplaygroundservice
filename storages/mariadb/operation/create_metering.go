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
	"gorm.io/gorm/clause"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

/*
CreateMetering 在 Metering 資料表寫入一筆 Metering Record

errors:
- 15000000(internal server error)
- 15000021(metering exist)
*/
func (o *Operation) CreateMetering(ctx context.Context, input *common.CreateMeteringInput) (output *common.CreateMeteringOutput, err error) {
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

	if createErr := o.conn.WithContext(ctx).Clauses(clause.Returning{}).Create(&input.Metering).Error; createErr != nil {
		if sqlErr, ok := createErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// Metering 已經存在
			case mariadb_com.ER_DUP_ENTRY:
				zap.L().With(
					zap.String(cnt.Storage, "o.conn.WithContext().Clauses(...).Create(...).Error"),
					zap.Any(cnt.RequestID, requestID),
					zap.String("name", GetTableName(&tables.Metering{})),
					zap.Any("value", input),
				).Warn(createErr.Error())
				err = tkErr.New(cnt.StorageMeteringExistErr).WithInner(createErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Clauses(...).Create(...).Error"),
			zap.Any(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Project{})),
			zap.Any("value", input),
		).Error(createErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(createErr)
		return
	}
	output = &common.CreateMeteringOutput{
		Metering: input.Metering,
	}
	return
}
