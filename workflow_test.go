package goauto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type noTask struct{}

func (t noTask) Run(*TaskInfo) (err error) {
	return
}

func TestNewWorkflow(t *testing.T) {
	wf := NewWorkflow(noTask{}, noTask{}, noTask{})
	err := wf.WatchPattern(".*", "/^[/rgsa*&")
	assert.NotNil(t, err)
	err = wf.WatchPattern(".*")
	assert.Nil(t, err)
	err = wf.WatchPattern(".*", ".*\\.go$")
	assert.Nil(t, err)
}
