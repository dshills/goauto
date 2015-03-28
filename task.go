// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

// A Task is an atomic unit of work
type Task struct {
	Banner     string                                       // Printed to wout before tasks starts
	Pass       string                                       // Printed to wout after task completes without error
	Fail       string                                       // Printed to wout after task completes with error
	Buffer     bytes.Buffer                                 // Output storage if any
	FileName   string                                       // Name of file after task is complete
	Target     string                                       // file name after running TargetFunc
	TaskFunc   func(task *Task, wout, werr io.Writer) error // The actual task
	TargetFunc func(string) string                          // Optional function to transform FileName to Target
}

// NewTaskType returns a new task using the provided task function
func NewTaskType(f func(t *Task, wout, werr io.Writer) error) *Task {
	nt := new(Task)
	nt.TaskFunc = f
	return nt
}

// Execute converts FileName to Target using the TargetFunc if provided, prints the Banner and runs the task function
func (t *Task) Execute(wout, werr io.Writer) error {
	t.Target = t.FileName
	if t.TargetFunc != nil {
		t.Target = t.TargetFunc(t.FileName)
	}
	if Verbose {
		log.Println("Execute Task", t)
	}
	if t.Banner != "" {
		fmt.Fprintln(wout, t.Banner, t.Target)
	}
	if err := t.TaskFunc(t, wout, werr); err != nil {
		if t.Fail != "" {
			fmt.Fprintln(wout, t.Fail)
		}
		return err
	}

	if t.Pass != "" {
		fmt.Fprintln(wout, t.Pass)
	}
	return nil
}
