package goauto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorkflow(t *testing.T) {
	assert.NotPanics(t, func() {
		NewWorkflow("WORKFLOW", ".*", &Task{})
	}, "Should not panic")

	assert.Panics(t, func() {
		NewWorkflow("WORKFLOW", "/^[/rgsa*&", &Task{})
	}, "Should panic")
}
