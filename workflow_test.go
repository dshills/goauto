package goauto

import "testing"

type noTask struct{}

func (t noTask) Run(*TaskInfo) (err error) {
	return
}

func TestNewWorkflow(t *testing.T) {
	wf := NewWorkflow(noTask{}, noTask{}, noTask{})
	err := wf.WatchPattern(".*", "/^[/rgsa*&")
	if err == nil {
		t.Errorf("Expected error for bad regexp\n")
	}
	err = wf.WatchPattern(".*")
	if err != nil {
		t.Error(err)
	}
	err = wf.WatchPattern(".*", ".*\\.go$")
	if err != nil {
		t.Error(err)
	}
}
