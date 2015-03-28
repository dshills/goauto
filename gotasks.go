// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"errors"
	"io"
	"os/exec"
	"strings"
)

// NewGoPrjTask returns a new task that will execute the given Go command
// The Target will be the relative directory from the $GOPATH
func NewGoPrjTask(cmd string) *Task {
	t := NewShellTask("go", cmd)
	t.TargetFunc = GoRelSrcDir
	t.Banner = strings.Title("go " + cmd + "...")
	return t
}

// NewGoTestTask returns a new task that will run all the project tests
func NewGoTestTask() *Task {
	return NewGoPrjTask("test")
}

// NewGoVetTask returns a new task that will vet the project
func NewGoVetTask() *Task {
	t := NewGoPrjTask("vet")
	t.Pass = "ok"
	return t
}

// NewGoBuildTask returns a task that will build the project
func NewGoBuildTask() *Task {
	t := NewGoPrjTask("build")
	t.Pass = "ok"
	return t
}

// NewGoInstallTask returns a task that will install the project
func NewGoInstallTask() *Task {
	t := NewGoPrjTask("install")
	t.Pass = "ok"
	return t
}

// NewGoLintTask returns a task that will golint the project
func NewGoLintTask() *Task {
	t := NewTaskType(func(t *Task, wout, werr io.Writer) error {
		t.Buffer.Reset()
		cmd := exec.Command("golint", t.Target)
		cmd.Stdout = &t.Buffer
		cmd.Stderr = werr
		defer func() {
			t.Buffer.WriteTo(wout)
		}()
		if err := cmd.Run(); err != nil {
			return err
		}
		if t.Buffer.Len() > 0 {
			return errors.New("FAIL")
		}
		return nil
	})
	t.TargetFunc = GoRelSrcDir
	t.Banner = "Go Lint..."
	t.Pass = "ok"
	return t
}

// NewGoFmtTask returns a task that will format the file
func NewGoFmtTask() *Task {
	t := NewShellTask("gofmt")
	t.Banner = "Go Format..."
	return t
}

// NewGoImportsTask returns a task that will insert imports into a file
func NewGoImportsTask() *Task {
	t := NewShellTask("goimports")
	t.Banner = "Go Imports..."
	return t
}
