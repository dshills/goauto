package goauto

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	p := NewPipeline("Pipline Name", Silent)
	assert.NotNil(t, p)
}

func TestPipelineRec(t *testing.T) {
	p := NewPipeline("Test Pipeline", Silent)
	tp := filepath.Join("src", "github.com", "dshills", "goauto")
	err := p.WatchRecursive(tp, IgnoreHidden)
	assert.Nil(t, err)
}

func TestPipelineWorkflow(t *testing.T) {
	wf := Workflow{}
	p := NewPipeline("Test Pipeline", Verbose)
	p.Add(&wf)

	wf2 := NewWorkflow(NewGoVetTask(), NewGoLintTask())
	p.Add(wf2)
}

func TestPipelineConcurrency(t *testing.T) {
	t0 := time.Now()
	p := NewPipeline("Test Pipeline", Silent)
	tp := filepath.Join("src", "github.com", "dshills", "goauto", "testing")
	err := p.WatchRecursive(tp, IgnoreHidden)
	assert.Nil(t, err)

	wf := NewWorkflow(NewGoVetTask(), NewGoLintTask(), NewGoBuildTask())
	p.Add(wf)

	p.Stop()

	go p.Start()

	atp, err := AbsPath(tp)
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		n := filepath.Join(atp, strconv.Itoa(i))
		os.Mkdir(n, 0744)
	}

	for i := 0; i < 100; i++ {
		n := filepath.Join(atp, strconv.Itoa(i))
		os.Remove(n)
	}

	p.Stop()

	go p.Start()

	for i := 0; i < 100; i++ {
		n := filepath.Join(atp, strconv.Itoa(i))
		os.Mkdir(n, 0744)
	}

	for i := 0; i < 100; i++ {
		n := filepath.Join(atp, strconv.Itoa(i))
		os.Remove(n)
	}

	time.Sleep(2 * time.Second)
	p.Stop()

	t1 := time.Now()
	log.Printf("TestPipelineConcurrency finished in %v", t1.Sub(t0))
}
