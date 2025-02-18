package fsm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	tests "pegasus-cloud.com/aes/toolkits/testing"
)

const (
	_name          = "test"
	_state1        = "s1"
	_state2        = "s2"
	_event1        = "e1"
	_empty         = ""
	_metadataKey   = "METADATA_KEY"
	_metadataValue = "METADATA_VALUE"
)

type mockFSM struct {
	tests.BaseSuite

	ctx context.Context
}

func TestCollectTemplateDataSuite(t *testing.T) {
	m := &mockFSM{}
	m.Origin = m
	suite.Run(t, m)
}

func (m *mockFSM) SetupTest() {
	m.ctx = context.Background()
}

func (m *mockFSM) TestName() {
	m.SetupTest()
	f := New(_name)
	m.Equal(_name, f.Name())
}

func (m *mockFSM) TestCurrent() {
	m.SetupTest()

	f := New(_name)
	f.NewState(_state1)
	f.NewState(_state2)
	f.NewEvent(_event1, _state1, _state2)

	f.Submit(_state1)
	m.Equal(_state1, f.Current())

	f.Event(m.ctx, _event1)
	m.Equal(_state2, f.Current())
}

func (m *mockFSM) TestIs() {
	m.SetupTest()

	f := New(_name)
	f.NewState(_state1)
	f.NewState(_state2)
	f.NewEvent(_event1, _state1, _state2)
	f.Submit(_state1)

	m.Equal(true, f.Is(_state1))
	m.Equal(false, f.Is(_state2))

	f.Event(m.ctx, _event1)
	m.Equal(false, f.Is(_state1))
	m.Equal(true, f.Is(_state2))
}

func (m *mockFSM) TestSetState() {
	m.SetupTest()

	f := New(_name)
	f.NewState(_state1)
	f.NewState(_state2)
	f.NewEvent(_event1, _state1, _state2)
	f.Submit(_state1)

	f.SetState(_state2)
	m.Equal(_state2, f.Current())
}

func (m *mockFSM) TestCan() {
	m.SetupTest()

	f := New(_name)
	f.NewState(_state1)
	f.NewState(_state2)
	f.NewEvent(_event1, _state1, _state2)
	f.Submit(_state1)

	m.Equal(true, f.Can(_event1))
}

func (m *mockFSM) TestMetadata() {
	m.SetupTest()

	f := New(_name)
	f.NewState(_state1)
	f.NewState(_state2)
	f.NewEvent(_event1, _state1, _state2)
	f.Submit(_state1)

	f.SetMetadata(_metadataKey, _metadataValue)

	val, exist := f.Metadata(_metadataKey)
	m.Equal(_metadataValue, val)
	m.Equal(true, exist)

	f.DeleteMetadata(_metadataKey)
	_, exist = f.Metadata(_metadataKey)
	m.Equal(false, exist)
}

func (m *mockFSM) TestStates() {
	m.SetupTest()

	f := New(_name)
	f.NewState(_state1)
	f.NewState(_state2)
	f.NewEvent(_event1, _state1, _state2)
	f.Submit(_state1)

	m.ElementsMatch([]string{_state1, _state2}, f.States())
}

func (m *mockFSM) TestEvents() {
	m.SetupTest()
	{
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		m.ElementsMatch([]string{_event1}, f.Events())
	}
	{
		tri := ""
		callback := func(msg string) func(ctx context.Context, e *Event) {
			return func(ctx context.Context, e *Event) {
				tri = msg
			}
		}

		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2).After(callback(_state1))
		f.NewEvent(_event1, _state2, _state1).After(callback(_state2))
		f.Submit(_state1)
		f.Event(m.ctx, _event1)

		m.Equal(_state1, tri)
		m.ElementsMatch([]string{_event1 + "_" + _state1, _event1 + "_" + _state2}, f.Events())
	}
}

func (m *mockFSM) TestCallbacks() {
	m.SetupTest()

	tri := []string{}

	callback := func(msg string) func(ctx context.Context, e *Event) {
		return func(ctx context.Context, e *Event) {
			tri = append(tri, msg)
		}
	}

	f := New(_name)
	f.EnterState(callback("enter_state")).LeaveState(callback("leave_state"))
	f.BeforeEvent(callback("before_event")).AfterEvent(callback("after_event"))

	f.NewState(_state1).Enter(callback("enter_" + _state1)).Leave(callback("leave_" + _state1))
	f.NewState(_state2).Enter(callback("enter_" + _state2)).Leave(callback("leave_" + _state2))
	f.NewEvent(_event1, _state1, _state2).Before(callback("before_" + _event1)).After(callback("after_" + _event1))

	f.Submit(_state1)
	f.Event(m.ctx, _event1)

	m.ElementsMatch(tri,
		[]string{
			"before_" + _event1,
			"before_event",
			"leave_" + _state1,
			"leave_state",
			"enter_" + _state2,
			"enter_state",
			"after_" + _event1,
			"after_event",
		})
}

func (m *mockFSM) TestVisualizeWithType() {
	m.SetupTest()

	f := New(_name)
	f.NewState(_state1)
	f.NewState(_state2)
	f.NewEvent(_event1, _state1, _state2)
	f.Submit(_state1)

	_, err := f.VisualizeWithType(GRAPHVIZ)
	m.Equal(nil, err)
}

func (m *mockFSM) TestError() {
	m.SetupTest()

	{ // Current()
		f := New(_name)
		m.Equal(_empty, f.Current())

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "Current()"}.Error())
	}
	{ // Is()
		f := New(_name)
		m.Equal(false, f.Is(_state1))

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "Is()"}.Error())
	}
	{ // SetState()
		f := New(_name)
		f.SetState(_state1)

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "SetState()"}.Error())
	}
	{ // Can()
		f := New(_name)
		f.Can(_event1)

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "Can()"}.Error())
	}
	{ // SetMetadata()
		f := New(_name)
		f.SetMetadata(_metadataKey, _metadataValue)

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "SetMetadata()"}.Error())
	}
	{ // Metadata()
		f := New(_name)
		f.Metadata(_metadataKey)

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "Metadata()"}.Error())
	}
	{ // DeleteMetadata()
		f := New(_name)
		f.DeleteMetadata(_metadataKey)

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "DeleteMetadata()"}.Error())
	}
	{ // States()
		f := New(_name)
		f.States()

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "States()"}.Error())
	}
	{ // Events()
		f := New(_name)
		f.Events()

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "Events()"}.Error())
	}
	{ // VisualizeWithType()
		f := New(_name)
		f.VisualizeWithType(GRAPHVIZ)

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "VisualizeWithType()"}.Error())
	}
	{ // NewEvent()
		f := New(_name)
		f.NewState(_state1)
		f.NewEvent(_event1, _state1, _state2)

		err := f.Submit(_state1)
		m.EqualError(err, InvalidStateError{State: _state2}.Error())
	}
	{ // NewEvent()
		f := New(_name)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)

		err := f.Submit(_state1)
		m.EqualError(err, InvalidStateError{State: _state1}.Error())
	}
	{ // LeaveState()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		f.LeaveState(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "LeaveState()"}.Error())
	}
	{ // EnterState()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		f.EnterState(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "EnterState()"}.Error())
	}
	{ // BeforeEvent()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		f.BeforeEvent(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "BeforeEvent()"}.Error())
	}
	{ // AfterEvent()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		f.AfterEvent(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "AfterEvent()"}.Error())
	}
	{ // NewState()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		f.NewState("")

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "NewState()"}.Error())
	}
	{ // NewEvent()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		f.NewEvent("", "", "")

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "NewEvent()"}.Error())
	}
	{ // Submit()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "Submit()"}.Error())
	}
	{ // Submit()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit("test")

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "initial state is not defined"}.Error())
	}
	{ // Event()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Event(m.ctx, _event1)

		err := f.Submit(_state1)
		m.EqualError(err, NotSubmittedError{Reason: "Event()"}.Error())
	}
	{ // Leave()
		f := New(_name)
		s := f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		s.Leave(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "Leave()"}.Error())
	}
	{ // Enter()
		f := New(_name)
		s := f.NewState(_state1)
		f.NewState(_state2)
		f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		s.Enter(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "Enter()"}.Error())
	}
	{ // Before()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		e := f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		e.Before(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "Before()"}.Error())
	}
	{ // AfterEvent()
		f := New(_name)
		f.NewState(_state1)
		f.NewState(_state2)
		e := f.NewEvent(_event1, _state1, _state2)
		f.Submit(_state1)

		e.After(nil)

		err := f.Submit(_state1)
		m.EqualError(err, SubmittedError{Reason: "After()"}.Error())
	}
}
