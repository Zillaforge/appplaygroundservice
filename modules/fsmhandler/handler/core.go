package handler

import (
	fsm "AppPlaygroundService/utility/fsm"
	"context"
)

type Method interface {
	// Name returns the name of the FSM.
	Name() string
	// States returns a list of all states
	States() []string
	// Events returns a list of all events
	Events() []string
	// VisualizeWithType outputs a visualization of a FSM in the desired format.
	VisualizeWithType() (string, error)

	// Event initiates a state transition with the named event.
	//
	// The call takes a variable number of arguments that will be passed to the
	// callback, if defined.
	//
	// It will return nil if the state change is ok or one of these errors:
	//
	// - event X inappropriate because previous transition did not complete
	//
	// - event X inappropriate in current state Y
	//
	// - event X does not exist
	//
	// - internal error on state transition
	//
	// The last error should never occur in this situation and is a sign of an
	// internal bug.
	Event(ctx context.Context, id, event string, args ...interface{}) error

	// Is returns true if state is the current state.
	Is(ctx context.Context, id, state string) (bool, error)

	// Can returns true if event can occur in the current state.
	Can(ctx context.Context, id, event string) bool

	// fsm constructs a FSM.
	FSM(initial string) *fsm.FSM
}
