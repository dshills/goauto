package goauto

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	p := NewPipeline("Pipline Name", Verbose)
	assert.NotNil(t, p)
}

func TestPipelineRec(t *testing.T) {
	p := NewPipeline("Test Pipeline", Verbose)
	tp := filepath.Join("src", "gituhub.com", "dshills", "goauto")
	p.WatchRecursive(tp, IgnoreHidden)
}

func TestPipelineWorkflow(t *testing.T) {
	wf := Workflow{}
	p := NewPipeline("Test Pipeline", Verbose)
	p.Add(&wf)

	wf2 := NewWorkflow(NewGoVetTask(), NewGoLintTask())
	p.Add(wf2)
}
