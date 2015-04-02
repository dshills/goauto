package goauto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	tp := filepath.Join("src", "github.com", "dshills", "goauto")
	p := NewPipeline("Pipline Name", tp, os.Stdout, os.Stderr, &Workflow{})
	assert.NotNil(t, p)
}

func TestPipelineRec(t *testing.T) {
	p := Pipeline{Name: "Test Pipeline"}
	tp := filepath.Join("src", "gituhub.com", "dshills", "goauto")
	p.WatchRecursive(tp, IgnoreHidden)
}

func TestPipelineWorkflow(t *testing.T) {
	wf := Workflow{}
	p := Pipeline{Name: "Test Pipeline"}
	p.AddWorkflow(&wf)

	wf2 := NewWorkflow(NewGoVetTask(), NewGoLintTask())
	p.AddWorkflow(wf2)
}
