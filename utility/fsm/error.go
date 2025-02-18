package fsm

import (
	"bytes"
	"strings"
)

type FSMErrors []error

func (e FSMErrors) Error() string {
	buff := bytes.NewBufferString("")

	for _, err := range e {
		buff.WriteString(err.Error())
		buff.WriteString("\n")
	}
	return strings.TrimSpace(buff.String())
}

// InvalidStateError is returned by FSM.NewEvent() when the event cannot fine in the defined state.
type InvalidStateError struct {
	State string
}

func (e InvalidStateError) Error() string {
	return "State(" + e.State + ") is not defined"
}

// InvalidEventError is returned by FSM.Event() when the event cannot be called
// in the current state.
type InvalidEventError struct {
	Event string
	State string
}

func (e InvalidEventError) Error() string {
	return "event " + e.Event + " inappropriate in current state " + e.State
}

// NotSubmittedError is returned when the fsm not yet submitted.
type NotSubmittedError struct {
	Reason string
}

func (e NotSubmittedError) Error() string {
	return e.Reason + ": not yet submitted"
}

// SubmittedError is returned when the fsm submitted.
type SubmittedError struct {
	Reason string
}

func (e SubmittedError) Error() string {
	return e.Reason + ": already submitted"
}
