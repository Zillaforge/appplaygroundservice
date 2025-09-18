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

// GetModuleCategory ...
func (o *Operation) GetModuleCategory(ctx context.Context, input *common.GetModuleCategoryInput) (output *common.GetModuleCategoryOutput, err error) {
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

	output = &common.GetModuleCategoryOutput{}
	if getErr := o.conn.WithContext(ctx).Model(&tables.ModuleCategory{}).Where("id = ?", input.ID).First(&output.ModuleCategory).Error; getErr != nil {
		if errors.Is(getErr, gorm.ErrRecordNotFound) {
			// ModuleCategory not found
			zap.L().With(
				zap.String(cnt.Storage, "record not found"),
				zap.String(cnt.RequestID, requestID),
				zap.String("name", GetTableName(&tables.ModuleCategory{})),
			).Error(getErr.Error())
			err = tkErr.New(cnt.StorageModuleCategoryNotFoundErr).WithInner(getErr)
			return
		}
		zap.L().With(
			zap.String(cnt.Storage, "o.conn.WithContext().Model(...).Where(...).First(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.ModuleCategory{})),
		).Error(getErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(getErr)
		return
	}
	return
}
