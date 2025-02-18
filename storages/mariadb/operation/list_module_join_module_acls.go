package operation

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// ListModuleJoinModuleAcls ...
func (o *Operation) ListModuleJoinModuleAcls(ctx context.Context, input *common.ListModuleJoinModuleAclsInput) (output *common.ListModuleJoinModuleAclsOutput, err error) {
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

	output = &common.ListModuleJoinModuleAclsOutput{}
	limit, offset := -1, 0
	if input.Pagination != nil {
		limit = input.Pagination.Limit
		offset = input.Pagination.Offset
	}

	tx := o.conn
	if input.ProjectID != nil {
		// Querying Module Allow for Specific Project
		tx = tx.Or("allow_project_id = ?", *input.ProjectID)
		// Querying Global Public Module
		tx = tx.Or("public = 1")
	}

	if listErr := whereCascade(o.conn.WithContext(ctx).Model(&tables.ModuleJoinModuleAcl{}).Where(tx), &input.Where).Count(&output.Count).Limit(limit).Offset(offset).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "module_created_at"}, Desc: true}).Find(&output.ModuleJoinModuleAcls).Error; listErr != nil {
		zap.L().With(
			zap.String(cnt.Storage, "whereCascade(o.conn.WithContext().Model(...).Where(...), ...).Count(...).Limit(...).Offset(...).Order(...).Find(...).Error"),
			zap.String(cnt.RequestID, requestID),
			zap.String("name", GetTableName(&tables.ModuleJoinModuleAcl{})),
			zap.Any("value", input),
		).Error(listErr.Error())
		err = tkErr.New(cnt.StorageInternalServerErr).WithInner(listErr)
		return
	}
	return
}
