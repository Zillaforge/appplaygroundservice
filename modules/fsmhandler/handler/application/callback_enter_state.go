package application

import (
	"AppPlaygroundService/utility/fsm"
	"context"

	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
)

func (h Handler) callBackEnterState() func(context.Context, *fsm.Event) {
	// update current state to db
	return func(ctx context.Context, e *fsm.Event) {
		var applicationID = _applicationID.Get(ctx)
		for _, arg := range e.Args {
			if val, ok := arg.(*pb.UpdateApplicationInput); ok {
				state := e.FSM.Current()
				val.ID = applicationID
				val.State = &state
				h.setDBData(&ctx, val)
				return
			}
		}

		h.setDBState(&ctx, e.FSM.Current())
	}
}
