package operation

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
)

// ListApplications ...
func (o *Operation) ListApplications(ctx context.Context, input *common.ListApplicationsInput) (output *common.ListApplicationsOutput, err error) {
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

	output = &common.ListApplicationsOutput{}
	limit, offset := -1, 0
	if input.Pagination != nil {
		limit = input.Pagination.Limit
		offset = input.Pagination.Offset
	}
	if listErr := whereCascade(o.conn.WithContext(ctx).Preload("Module").Preload("Module.ModuleCategory").Model(&tables.Application{}), &input.Where).Count(&output.Count).Limit(limit).Offset(offset).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&output.Applications).Error; listErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Preload(...).Preload(...).Model(...), ...).Count(...).Limit(...).Offset(...).Order(...).Find(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.Application{})),
			zap.Any("value", input),
		).Error(listErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(listErr)
		return
	}
	return
}
