package fsm

import (
	"context"

	"github.com/looplab/fsm"
)

type Event fsm.Event
type Callback func(context.Context, *Event)

func (f Callback) Callback() func(context.Context, *fsm.Event) {
	return func(ctx context.Context, e *fsm.Event) {
		f(ctx, (*Event)(e))
	}
}

type VisualizeType string

const (
	// GRAPHVIZ the type for graphviz output (http://www.webgraphviz.com/)
	GRAPHVIZ VisualizeType = "graphviz"
	// MERMAID the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	MERMAID VisualizeType = "mermaid"
	// MermaidStateDiagram the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	MermaidStateDiagram VisualizeType = "mermaid-state-diagram"
	// MermaidFlowChart the type for mermaid output (https://mermaid-js.github.io/mermaid/#/flowchart) in the flow chart form
	MermaidFlowChart VisualizeType = "mermaid-flow-chart"
)

type (
	FSM struct {
		// name of FSM
		name string
		// state record all of states
		state map[string]*state
		// event record all of events
		event map[string]*event
		// callback record global callbacks
		callbacks map[string]Callback
		// submit is flag for fsm summit or not
		submit bool
		// err record all of error
		err FSMErrors
		// fsm is core for all feature
		fsm *fsm.FSM
	}
	iFSM interface {
		Name() string
		States() []string
		Events() []string
		VisualizeWithType(visualizeType VisualizeType) (string, error)

		NewState(name string) iState
		NewEvent(name, src, dst string) iEvent

		LeaveState(Callback) iFSM
		EnterState(Callback) iFSM
		BeforeEvent(Callback) iFSM
		AfterEvent(Callback) iFSM

		Submit(initialState string) error
		Event(ctx context.Context, event string, args ...interface{}) error

		Current() string
		Is(state string) bool
		Can(event string) bool
		SetState(state string)

		Metadata(key string) (interface{}, bool)
		SetMetadata(key string, dataValue interface{})
		DeleteMetadata(key string)
	}
)

type (
	state struct {
		fsm *FSM

		// name of state
		name string
		// // enter called after entering <STATE>
		enter Callback
		// leave called before leaving <STATE>
		leave Callback
	}
	iState interface {
		Enter(Callback) iState
		Leave(Callback) iState
	}
)

type (
	event struct {
		fsm *FSM

		// name of state
		name string
		// src is a slice of source states
		src string
		// dst is the destination state
		dst string
		// before called before event named <EVENT>
		before Callback
		// after called after event named <EVENT>
		after Callback

		duplicated bool
	}
	iEvent interface {
		Before(Callback) iEvent
		After(Callback) iEvent
	}
)

const (
	leaveState  = "leave_"
	enterState  = "enter_"
	beforeEvent = "before_"
	afterEvent  = "after_"
)

// NewFSM constructs a FSM.
func New(name string) *FSM {
	return &FSM{
		name:      name,
		state:     map[string]*state{},
		event:     map[string]*event{},
		callbacks: map[string]Callback{},
	}
}

// Name returns the name of the FSM.
func (f FSM) Name() string {
	return f.name
}

// Current returns the current state of the FSM.
func (f *FSM) Current() string {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "Current()"})
		return ""
	}
	return f.fsm.Current()
}

// Is returns true if state is the current state.
func (f *FSM) Is(state string) bool {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "Is()"})
		return false
	}
	return f.fsm.Is(state)
}

// SetState allows the user to move to the given state from current state.
// The call does not trigger any callbacks, if defined.
func (f *FSM) SetState(state string) {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "SetState()"})
		return
	}
	f.fsm.SetState(state)
}

// Can returns true if event can occur in the current state.
func (f *FSM) Can(event string) bool {
	if !f.submit {
		err := NotSubmittedError{Reason: "Can()"}
		f.err = append(f.err, err)
		return false
	}

	if val, ok := f.event[event]; !ok {
		err := InvalidEventError{Event: event, State: f.fsm.Current()}
		f.err = append(f.err, err)
		return false
	} else {
		if val.duplicated {
			name := event + "_" + f.fsm.Current()
			if _, exist := f.event[name]; exist {
				event = name
			}
		}
	}

	return f.fsm.Can(event)
}

// Metadata returns the value stored in metadata
func (f *FSM) Metadata(key string) (interface{}, bool) {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "Metadata()"})
		return nil, false
	}
	return f.fsm.Metadata(key)
}

// SetMetadata stores the dataValue in metadata indexing it with key
func (f *FSM) SetMetadata(key string, dataValue interface{}) {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "SetMetadata()"})
		return
	}
	f.fsm.SetMetadata(key, dataValue)
}

// DeleteMetadata deletes the dataValue in metadata by key
func (f *FSM) DeleteMetadata(key string) {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "DeleteMetadata()"})
		return
	}
	f.fsm.DeleteMetadata(key)
}

// States returns a list of all states
func (f *FSM) States() (states []string) {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "States()"})
		return
	}

	for state := range f.state {
		states = append(states, state)
	}
	return
}

// Events returns a list of all events
func (f *FSM) Events() (events []string) {
	if !f.submit {
		f.err = append(f.err, NotSubmittedError{Reason: "Events()"})
		return
	}

	for name, event := range f.event {
		if !event.duplicated {
			events = append(events, name)
		}
	}
	return
}

// VisualizeWithType outputs a visualization of a FSM in the desired format.
// If the type is not given it defaults to GRAPHVIZ
func (f *FSM) VisualizeWithType(visualizeType VisualizeType) (string, error) {
	if !f.submit {
		err := NotSubmittedError{Reason: "VisualizeWithType()"}
		f.err = append(f.err, err)
		return "", err
	}

	return fsm.VisualizeWithType(f.fsm, fsm.VisualizeType(visualizeType))
}

// LeaveState called before leaving all states
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (f *FSM) LeaveState(callback Callback) iFSM {
	if f.submit {
		f.err = append(f.err, SubmittedError{Reason: "LeaveState()"})
		return f
	}

	if callback != nil {
		f.callbacks[leaveState+"state"] = callback
	}
	return f
}

// EnterState called after entering all states
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (f *FSM) EnterState(callback Callback) iFSM {
	if f.submit {
		f.err = append(f.err, SubmittedError{Reason: "EnterState()"})
		return f
	}

	if callback != nil {
		f.callbacks[enterState+"state"] = callback
	}
	return f
}

// BeforeEvent called before all events
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (f *FSM) BeforeEvent(callback Callback) iFSM {
	if f.submit {
		f.err = append(f.err, SubmittedError{Reason: "BeforeEvent()"})
		return f
	}

	if callback != nil {
		f.callbacks[beforeEvent+"event"] = callback
	}
	return f
}

// AfterEvent called after all events
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (f *FSM) AfterEvent(callback Callback) iFSM {
	if f.submit {
		f.err = append(f.err, SubmittedError{Reason: "AfterEvent()"})
		return f
	}

	if callback != nil {
		f.callbacks[afterEvent+"event"] = callback
	}
	return f
}

// NewState constructs a state from name and callback
func (f *FSM) NewState(name string) iState {
	if f.submit {
		f.err = append(f.err, SubmittedError{Reason: "NewState()"})
		return nil
	}

	f.state[name] = &state{fsm: f, name: name}
	return f.state[name]
}

// NewEvent constructs a event from name, src_state, dst_state and callback
//
// Rename <EVENT> to <EVENT_SRC> when event name duplicated
func (f *FSM) NewEvent(name, src, dst string) iEvent {
	if f.submit {
		f.err = append(f.err, SubmittedError{Reason: "NewEvent()"})
		return nil
	}

	if _, ok := f.state[src]; !ok {
		f.err = append(f.err, InvalidStateError{State: src})
	}
	if _, ok := f.state[dst]; !ok {
		f.err = append(f.err, InvalidStateError{State: dst})
	}

	if val, exist := f.event[name]; exist {
		if !val.duplicated {
			f.event[name] = &event{fsm: f, name: name, duplicated: true}

			name := val.name + "_" + val.src
			f.event[name] = val
			f.event[name].name = name
		}
		name = name + "_" + src
	}
	f.event[name] = &event{fsm: f, name: name, src: src, dst: dst}

	return f.event[name]
}

// Submit constructs a FSM from states, events and callbacks.
func (f *FSM) Submit(initial string) error {
	if f.err != nil {
		return f.err
	}

	if f.submit {
		err := SubmittedError{Reason: "Submit()"}
		f.err = append(f.err, err)
		return err
	}

	callbacks := map[string]fsm.Callback{}
	events := []fsm.EventDesc{}
	initialStateRegister := false

	for _, state := range f.state {
		if state.name == initial {
			initialStateRegister = true
		}
		if state.enter != nil {
			callbacks[enterState+state.name] = state.enter.Callback()
		}
		if state.leave != nil {
			callbacks[leaveState+state.name] = state.leave.Callback()
		}
	}
	if !initialStateRegister {
		err := SubmittedError{Reason: "initial state is not defined"}
		f.err = append(f.err, err)
		return err
	}

	for _, event := range f.event {
		if event.duplicated {
			continue
		}
		events = append(events, fsm.EventDesc{
			Name: event.name,
			Src:  []string{event.src},
			Dst:  event.dst,
		})
		if event.before != nil {
			callbacks[beforeEvent+event.name] = event.before.Callback()
		}
		if event.after != nil {
			callbacks[afterEvent+event.name] = event.after.Callback()
		}
	}
	for name, callback := range f.callbacks {
		callbacks[name] = callback.Callback()
	}

	f.fsm = fsm.NewFSM(initial, events, callbacks)
	f.submit = true
	return nil
}

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
func (f *FSM) Event(ctx context.Context, event string, args ...interface{}) error {
	if !f.submit {
		err := NotSubmittedError{Reason: "Event()"}
		f.err = append(f.err, err)
		return err
	}

	if val, ok := f.event[event]; !ok {
		err := InvalidEventError{Event: event, State: f.fsm.Current()}
		f.err = append(f.err, err)
		return err
	} else {
		if val.duplicated {
			name := event + "_" + f.fsm.Current()
			if _, exist := f.event[name]; exist {
				event = name
			}
		}
	}

	return f.fsm.Event(ctx, event, args...)
}

// Leave called before leaving <OLD_STATE>
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (s *state) Leave(callback Callback) iState {
	if s.fsm.submit {
		s.fsm.err = append(s.fsm.err, SubmittedError{Reason: "Leave()"})
		return nil
	}

	s.leave = callback
	return s
}

// Enter called after entering <NEW_STATE>
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (s *state) Enter(callback Callback) iState {
	if s.fsm.submit {
		s.fsm.err = append(s.fsm.err, SubmittedError{Reason: "Enter()"})
		return nil
	}

	s.enter = callback
	return s
}

// Before called before event named <EVENT>
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (e *event) Before(callback Callback) iEvent {
	if e.fsm.submit {
		e.fsm.err = append(e.fsm.err, SubmittedError{Reason: "Before()"})
		return nil
	}

	e.before = callback
	return e
}

// After called after event named <EVENT>
//
// Callbacks are added as a map specified as Callbacks where the key is parsed
// as the callback event as follows, and called in the same order:
//
// 1. before_<EVENT> - called before event named <EVENT>
//
// 2. before_event - called before all events
//
// 3. leave_<OLD_STATE> - called before leaving <OLD_STATE>
//
// 4. leave_state - called before leaving all states
//
// 5. enter_<NEW_STATE> - called after entering <NEW_STATE>
//
// 6. enter_state - called after entering all states
//
// 7. after_<EVENT> - called after event named <EVENT>
//
// 8. after_event - called after all events
func (e *event) After(callback Callback) iEvent {
	if e.fsm.submit {
		e.fsm.err = append(e.fsm.err, SubmittedError{Reason: "After()"})
		return nil
	}

	e.after = callback
	return e
}
