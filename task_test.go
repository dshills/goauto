package goauto

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskType(t *testing.T) {
	nt := NewTaskType(func(t *Task, wout, werr io.Writer) error { return nil })
	err := nt.TaskFunc(nt, os.Stdout, os.Stderr)
	assert.Nil(t, err, "Expecting no error from TaskFunc")
}

func TestExecute(t *testing.T) {
	tf := func(t *Task, wout, werr io.Writer) error {
		t.Buffer.Reset()
		t.Buffer.WriteString(t.Target)
		t.Buffer.WriteString(t.Banner)
		t.Buffer.WriteString(t.FileName)
		fmt.Fprintln(wout, t.Target)
		fmt.Fprintln(werr, t.Target)
		return nil
	}
	nt := NewTaskType(tf)
	nt.Banner = "Test Task"
	err := nt.Execute(os.Stdout, os.Stderr)
	assert.Nil(t, err, "Expecting no error from Execute")
	assert.Equal(t, nt.Buffer.String(), nt.Target+nt.Banner+nt.FileName)
}

func TestTarget(t *testing.T) {
	tf := func(t *Task, wout, werr io.Writer) error {
		t.Buffer.Reset()
		t.Buffer.WriteString(t.Target)
		return nil
	}
	nt := NewTaskType(tf)
	nt.TargetFunc = strings.Title
	nt.FileName = "this is a test for the target function"
	err := nt.Execute(os.Stdout, os.Stderr)
	assert.Nil(t, err, "Expecting no error from Execute")
	assert.Equal(t, nt.Buffer.String(), nt.TargetFunc(nt.FileName))
}
