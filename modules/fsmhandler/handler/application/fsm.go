package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/fsmhandler"
	applicationCom "AppPlaygroundService/modules/fsmhandler/common/application"
	fsm "AppPlaygroundService/utility/fsm"
	"context"

	"github.com/Zillaforge/appplaygroundserviceclient/aps"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

type (
	ctxApplication string

	Handler struct {
		initial string
	}
)

const (
	_userID          fsmhandler.CtxID = "ctxUserID"
	_applicationID   fsmhandler.CtxID = "ctxApplicationID"
	_applicationData ctxApplication   = "ctxApplicationData"
)

func Init() {
	fsmhandler.Application = New(applicationCom.ReviewState)
}

func New(initial string) *Handler {
	return &Handler{
		initial: initial,
	}
}

func (h *Handler) FSM(initial string) *fsm.FSM {
	f := fsm.New(cnt.Application)

	f.NewState(applicationCom.ReviewState)
	f.NewState(applicationCom.RejectState)
	f.NewState(applicationCom.ProcessState)
	f.NewState(applicationCom.ReadyState)
	f.NewState(applicationCom.FailState).Enter(h.callBackEnterFailState())

	f.EnterState(h.callBackEnterState())

	f.NewEvent(applicationCom.ApproveEvent, applicationCom.ReviewState, applicationCom.ProcessState).After(h.callBackAfterApproveEvent())
	f.NewEvent(applicationCom.RejectEvent, applicationCom.ReviewState, applicationCom.RejectState)
	f.NewEvent(applicationCom.FinishEvent, applicationCom.ProcessState, applicationCom.ReadyState)
	f.NewEvent(applicationCom.FailEvent, applicationCom.ProcessState, applicationCom.FailState)

	if err := f.Submit(initial); err != nil {
		panic(err.Error())
	}

	return f
}

func (h *Handler) Name() string {
	return h.FSM(h.initial).Name()
}

func (h *Handler) States() []string {
	return h.FSM(h.initial).States()
}

func (h *Handler) Events() []string {
	return h.FSM(h.initial).Events()
}

func (h *Handler) VisualizeWithType() (string, error) {
	return h.FSM(h.initial).VisualizeWithType(fsm.MERMAID)
}

func (h *Handler) Event(ctx context.Context, id, event string, args ...interface{}) error {
	ctx = context.WithValue(ctx, _applicationID, id)

	state, err := h.getDBState(&ctx)
	if err != nil {
		return err
	}
	return h.FSM(state).Event(ctx, event, args...)
}

func (h *Handler) Is(ctx context.Context, id, targetState string) (bool, error) {
	ctx = context.WithValue(ctx, _applicationID, id)

	state, err := h.getDBState(&ctx)
	if err != nil {
		return false, err
	}
	return h.FSM(state).Is(targetState), nil
}

func (h *Handler) Can(ctx context.Context, id, event string) bool {
	ctx = context.WithValue(ctx, _applicationID, id)

	state, err := h.getDBState(&ctx)
	if err != nil {
		return false
	}
	return h.FSM(state).Can(event)
}

func (h *Handler) getDBState(ctx *context.Context) (string, error) {
	input := &pb.GetWithLanguageInput{ID: _applicationID.Get(*ctx)}
	output, err := aps.GetApplication(input, *ctx)
	if err != nil {
		return "", err
	}
	*ctx = context.WithValue(*ctx, _applicationData, output)
	return output.Application.State, nil
}

func (h *Handler) setDBState(ctx *context.Context, state string) error {
	input := &pb.UpdateApplicationInput{ID: _applicationID.Get(*ctx), State: &state, UpdaterID: _userID.Get(*ctx)}
	if _, err := h.setDBData(ctx, input); err != nil {
		return err
	}
	return nil
}

func (h *Handler) setDBData(ctx *context.Context, data *pb.UpdateApplicationInput) (*pb.ApplicationDetail, error) {
	if data.UpdaterID == "" {
		c := *ctx
		data.UpdaterID = c.Value(_applicationData).(*pb.ApplicationDetail).Application.UpdaterID
	}
	output, err := aps.UpdateApplication(data, *ctx)
	if err != nil {
		return nil, err
	}
	*ctx = context.WithValue(*ctx, _applicationData, output)
	return output, nil
}
