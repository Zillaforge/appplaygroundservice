package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/fsm"
	"context"

	"go.uber.org/zap"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (h Handler) callBackEnterFailState() func(context.Context, *fsm.Event) {
	// update current state to db
	funcName := tkUtils.NameOfFunction().Name()
	return func(ctx context.Context, e *fsm.Event) {
		var (
			requestID = utility.MustGetContextRequestID(ctx)
			applicationID  = _applicationID.Get(ctx)
		)

		zap.L().With(
			zap.String(cnt.FSM, funcName),
			zap.Any(cnt.RequestID, requestID),
			zap.Any("application-id", applicationID),
		).Info("application fail")
	}
}
