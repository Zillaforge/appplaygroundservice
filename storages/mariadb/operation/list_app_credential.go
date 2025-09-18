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

// ListAppCredentials ...
func (o *Operation) ListAppCredentials(ctx context.Context, input *common.ListAppCredentialsInput) (output *common.ListAppCredentialsOutput, err error) {
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

	output = &common.ListAppCredentialsOutput{}
	limit, offset := -1, 0
	if input.Pagination != nil {
		limit = input.Pagination.Limit
		offset = input.Pagination.Offset
	}
	if listErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.AppCredential{}), &input.Where).Count(&output.Count).Limit(limit).Offset(offset).Find(&output.AppCredentials).Error; listErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...), ...).Count(...).Limit(...).Offset(...).Find(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.AppCredential{})),
			zap.Any("value", input),
		).Error(listErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(listErr)
		return
	}
	return
}
