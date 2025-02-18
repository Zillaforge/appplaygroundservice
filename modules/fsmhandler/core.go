package fsmhandler

import (
	"AppPlaygroundService/modules/fsmhandler/handler"
	"context"
)

var (
	Application handler.Method
)

type CtxID string

func (c CtxID) Get(ctx context.Context) string {
	v := ctx.Value(c)
	if v == nil {
		return ""
	}
	return v.(string)
}
