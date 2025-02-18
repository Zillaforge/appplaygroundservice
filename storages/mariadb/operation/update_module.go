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

// UpdateModule ...
func (o *Operation) UpdateModule(ctx context.Context, input *common.UpdateModuleInput) (output *common.UpdateModuleOutput, err error) {
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

	whereCondition := &common.UpdateModuleInput{
		ID: input.ID,
	}
	output = &common.UpdateModuleOutput{}
	if updateErr := o.conn.WithContext(ctx).Model(&tables.Module{}).Where(queryConversion(*whereCondition)).Updates(queryConversion(*input.UpdateData)).Preload("ModuleCategory").First(&output.Module).Error; updateErr != nil {
		if sqlErr, ok := updateErr.(*mysql.MySQLError); ok {
			switch sqlErr.Number {
			// module 已經存在
			case mariadb_com.ER_DUP_ENTRY:
				zap.L().With(
					zap.String(cnt.Storage, " o.conn.WithContext().Model(...).Where(...).Updates(...).First(...).Error"),
					zap.String("name", GetTableName(&tables.Module{})),
					zap.Any("value", input),
					zap.Any(cnt.RequestID, requestID),
				).Warn(updateErr.Error())
				err = tkErr.New(cnt.StorageModuleExistErr).WithInner(updateErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).Updates(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Module{})),
			zap.Any("value", input),
		).Error(updateErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(updateErr)
		return
	}
	return
}
