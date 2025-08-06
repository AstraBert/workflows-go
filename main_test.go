package workflowsgo

import (
	"maps"
	"slices"
	"testing"
)

func TestEvent(t *testing.T) {
	var tests = []struct {
		firstStep string
		data      map[string]string
		want      []string
	}{
		{"step1", map[string]string{"key": "hello"}, []string{"step1", "hello"}},
		{"step2", map[string]string{"key": "ciao"}, []string{"step2", "ciao"}},
		{"step3", map[string]string{"key": "guten tag"}, []string{"step3", "guten tag"}},
	}

	for _, tt := range tests {
		event := NewBaseEvent(tt.firstStep, tt.data)
		_, err := any(*event).(GenericEvent)
		if err {
			t.Error("Testing NewBaseEvent: NewBaseEvent does not return an instance of GenericEvent")
		}
		stepName := event.NextStep
		keyVal, _ := event.Get("key")
		if v, ok := keyVal.(string); !ok {
			t.Errorf("Testing BaseEvent.Get: Expecting the returned value to be a string.")
		} else {
			if !slices.Equal(tt.want, []string{stepName, v}) {
				t.Errorf("Testing NewBaseEvent and BaseEvent.Get: want %v, got %v", tt.want, []string{stepName, v})
			}
		}
	}
}

type User struct {
	name  string
	email string
	age   int
}

func TestContext(t *testing.T) {
	state := map[string]any{
		"nMessages":  4000,
		"success":    true,
		"iterations": 3,
	}
	store := map[string]any{
		"user": User{"John Doe", "john.doe@example.com", 30},
	}
	ctx := NewBaseContext(store, state)
	_, err := any(*ctx).(GenericContext)
	if err {
		t.Error("Testing NewBaseContext: NewBaseContext does not return an instance of GenericContext")
	}
	if !maps.Equal(state, ctx.GetState()) {
		t.Errorf("Testing NewBaseContext and BaseContext.GetState: got %v, want %v", ctx.GetState(), state)
	}
	modifiedState := map[string]any{
		"nMessages":  3500,
		"success":    true,
		"iterations": 4,
	}
	ctx.SetState(modifiedState)
	if !maps.Equal(modifiedState, ctx.GetState()) {
		t.Errorf("Testing BaseContext.SetState: got %v, want %v", ctx.GetState(), modifiedState)
	}
	user, _ := ctx.GetValue("user")
	if user != store["user"] {
		t.Errorf("Testing BaseContext.GetValue: got %v, want %v", user, store["user"])
	}
	modifiedUser := User{"Johanna Doe", "johanna.doe@example.com", 25}
	ctx.StoreValue("user", modifiedUser)
	user1, _ := ctx.GetValue("user")
	if user1 != modifiedUser {
		t.Errorf("Testing BaseContext.StoreValue (modify existing value): got %v, want %v", user1, modifiedUser)
	}
	ctx.StoreValue("apiCredit", 200)
	apiCredit, _ := ctx.GetValue("apiCredit")
	if apiCredit != 200 {
		t.Errorf("Testing BaseContext.StoreValue (create a new value): got %d, want %d", apiCredit, 200)
	}
}

func mockStep(ev *BaseEvent, ctx *BaseContext) *BaseEvent {
	return NewBaseEvent("end", map[string]string{"output": "hello world"})
}

func TestWorkflow(t *testing.T) {
	wf := NewBaseWorkflow("firstStep", NewBaseContext(map[string]any{}, map[string]any{}), map[string]func(*BaseEvent, *BaseContext) *BaseEvent{"firstStep": mockStep})
	_, err := any(wf).(BaseWorkflow)
	if err {
		t.Error("Testing NewBaseWorkflow: NewBaseWorkflow does not return an instance of BaseWorkflow")
	}
	_, ok := any(wf).(GenericWorkflow)
	if ok {
		t.Error("Testing NewBaseWorkflow: NewBaseWorkflow does not return an instance of GenericWorkflow")
	}
	valid, _ := wf.Validate()
	if !valid {
		t.Error("Testing BaseWorkflow.Validate: BaseWorkflow is not valid, but it should be")
	}
	event := wf.TakeStep("firstStep", NewBaseEvent("mockEvent", map[string]string{"mock": "event"}), NewBaseContext(map[string]any{}, map[string]any{}))
	if val, _ := event.Get("output"); val != "hello world" {
		t.Errorf("Testing BaseWorkflow.TakeStep: Expected 'hello world', gotten %s", val)
	}
	if event.NextStep != "end" {
		t.Errorf("Testing BaseWorkflow.TakeStep: Expected 'end', gotten %s", event.NextStep)
	}
	output := wf.Output(event, NewBaseContext(map[string]any{}, map[string]any{}))
	if output != "hello world" {
		t.Errorf("Testing BaseWorkflow.Output: Expected 'hello world', gotten %s", output)
	}
	startCallBacks := []string{}
	endCallBacks := []string{}
	outputCallBacks := []any{}

	startEventCallBack := func(ev *BaseEvent) {
		startCallBacks = append(startCallBacks, ev.NextStep)
	}

	endEventCallBack := func(ev *BaseEvent) {
		endCallBacks = append(endCallBacks, event.NextStep)
	}

	outputCallBack := func(out any) {
		outputCallBacks = append(outputCallBacks, out)
	}

	wf.Run(NewBaseEvent("mockEvent", map[string]string{"mock": "event"}), NewBaseContext(map[string]any{}, map[string]any{}), startEventCallBack, endEventCallBack, outputCallBack)
	if !slices.Equal(startCallBacks, []string{"end"}) || !slices.Equal(outputCallBacks, []any{"hello world"}) || len(endCallBacks) != 0 {
		t.Errorf("Testing for BaseWorkflow.Run: want %v, %v, %d\ngot %v, %v, %d", []string{"end"}, []string{"hello world"}, 0, startCallBacks, outputCallBacks, len(endCallBacks))
	}
}
