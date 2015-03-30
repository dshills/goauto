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
	assert.NotPanics(t, func() {
		NewWorkflow("WORKFLOW", ".*", noTask{})
	}, "Should not panic")

	assert.Panics(t, func() {
		NewWorkflow("WORKFLOW", "/^[/rgsa*&", noTask{})
	}, "Should panic")
}
