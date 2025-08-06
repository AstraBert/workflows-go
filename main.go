// workflows-go is a package intended to build event-driven workflows in Go.
//
// The main use case for workflows-go is the development of AI-powered
// workflows, but the concept of stepwise and concatenated function
// execution can be applied to many fields of software engineering.
package workflowsgo

import (
	"errors"
	"fmt"
)

// GenericEvent is an interface that must be implemented by all structs
// representing events if they want to be considered as correct
// input/output for a workflow.
type GenericEvent interface {
	// Get is a function that fetches data stored within an Event, and returns that data. Data should be stored with strings as keys.
	Get(string) (any, bool)
}

// BaseEvent is a base implementation of the GenericEvent: it comes with NextStep (string) and Data (map) as attributes where the key information is stored.
type BaseEvent struct {
	NextStep string
	Data     map[string]string
}

// Get is a method of BaseEvent that fetches data stored within an BaseEvent.Data, and returns that data.
func (ev *BaseEvent) Get(key string) (val any, success bool) {
	val, success = ev.Data[key]
	return
}

// NewBaseEvent is a constructor function that, given a string representing the name of the following step and a map containing data, returns a BaseEvent.
func NewBaseEvent(followingStep string, data map[string]string) *BaseEvent {
	return &BaseEvent{
		Data:     data,
		NextStep: followingStep,
	}
}

// GenericContext is the interface representing a context, i.e. a storage
// space that is aimed at allowing persistency and stefulness for
// workflow executions. All structs representing a workflow context
// should implement this interface.
type GenericContext interface {
	// StoreValue is aimed at persistant storage of a key-value pair.
	StoreValue(string, any)

	// GetValue retrieves the value associated with a key within the persistent
	// storage
	GetValue(string) (any, bool)

	// GetState retrieves the in-BaseContext state of the workflow run.
	GetState() map[string]any

	// SetState allows to set a new state within the BaseContext.
	SetState(map[string]any)
}

// BaseContext is a base implementation of the GenericContext implementation.
// It comes with a Store (persistent storage), which is a map
// and a State (stateful execution) which also is a map.
type BaseContext struct {
	Store map[string]any
	State map[string]any
}

// NewBaseContext is a constructor that, starting from a map representing
// the Store and one representing the State, returns a BaseContext.
func NewBaseContext(store, state map[string]any) *BaseContext {
	return &BaseContext{
		Store: store,
		State: state,
	}
}

// StoreValue stores a key-value pair in BaseContext.Store.
func (ctx *BaseContext) StoreValue(key string, val any) {
	ctx.Store[key] = val
}

// GetValue fetches the value associated with a key in BaseContext.Store.
func (ctx *BaseContext) GetValue(key string) (val any, success bool) {
	val, success = ctx.Store[key]
	return
}

// GetState fetches BaseContext.State.
func (ctx *BaseContext) GetState() map[string]any {
	return ctx.State
}

// SetState assigns a value to BaseContext.State.
func (ctx *BaseContext) SetState(state map[string]any) {
	ctx.State = state
}

// GenericWorkflow is an interface providing a generic implementation of
// an event-driven workflow. Every struct representing a workflow should
// implement the GenericWorkflow interface.
//
// A workflow should be composed of steps, functions that take a GenericEvent and GenericContext as arguments: these steps should
// be in some way associated with strings representing their names.
type GenericWorkflow interface {
	// TakeStep allows separate execution single steps. It executes a step by selecting it with its name and passing a GenericEvent and a GenericContext to it.
	TakeStep(string, *GenericEvent, *GenericContext) *BaseEvent

	// Validate ensures that the structure of the workflow is correct
	Validate() bool

	// Run runs the workflow until completion. It takes an input event and an
	//  initial context, as well as three callback function, respectively
	//  for when an event starts being processed, for when a new event is
	//  emitted and for the workflow output
	Run(*GenericEvent, *GenericContext, func(*GenericEvent), func(*GenericEvent), func(any))

	// Output runs at the end of the workflow and returns the actual
	//  workflow output.
	Output(*GenericEvent, *GenericContext) any
}

// BaseWorkflow offers a base implementation of GenericWorkflow.
type BaseWorkflow struct {
	FirstStep string
	Context   *BaseContext
	Steps     map[string]func(*BaseEvent, *BaseContext) *BaseEvent
}

// Validate checks that the steps in the workflow are not named with 'end',
// a keyword reserved for the name of the output step.
func (wf *BaseWorkflow) Validate() (bool, error) {
	for k := range wf.Steps {
		if k == "end" {
			return false, errors.New("`end` is a reserved keyword, you cannot use it as a name for your steps")
		}
	}
	return true, nil
}

// TakeStep allows separate execution single steps by calling
// them with their name.
func (wf *BaseWorkflow) TakeStep(stepName string, ev *BaseEvent, ctx *BaseContext) *BaseEvent {
	step, ok := wf.Steps[stepName]
	if !ok {
		var data map[string]string = map[string]string{
			"output": fmt.Sprintf("There was an error while executing step %s: the step does not exist", stepName),
		}
		return NewBaseEvent("end", data)
	} else {
		return step(ev, ctx)
	}
}

// Run runs the workflow through completion
func (wf *BaseWorkflow) Run(inputEvent *BaseEvent, context *BaseContext, onEventStartCallBack func(*BaseEvent), onEventEndCallBack func(*BaseEvent), onOutputCallBack func(any)) {
	event := wf.TakeStep(wf.FirstStep, inputEvent, context)
	for {
		onEventStartCallBack(event)
		if event.NextStep == "end" {
			output := wf.Output(event, context)
			onOutputCallBack(output)
			break
		}
		event = wf.TakeStep(event.NextStep, event, context)
		onEventEndCallBack(event)
	}
}

// Output produces the output of the workflow.
func (wf *BaseWorkflow) Output(ev *BaseEvent, ctx *BaseContext) any {
	if ev.NextStep == "end" {
		output, ok := ev.Get("output")
		if ok {
			return output
		}
		return "No output produced"
	}
	return "Not an output step"
}

func NewBaseWorkflow(firstStep string, ctx *BaseContext, steps map[string]func(*BaseEvent, *BaseContext) *BaseEvent) *BaseWorkflow {
	return &BaseWorkflow{
		FirstStep: firstStep,
		Context:   ctx,
		Steps:     steps,
	}
}
