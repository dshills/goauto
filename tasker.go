// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"bytes"
	"io"
)

// A TaskInfo contains the results of running a Task
type TaskInfo struct {
	Src        string
	Target     string
	Buf        bytes.Buffer
	Tout, Terr io.Writer
}

// A Runner represents the function needed to satisfy a Tasker interface
type Runner func(*TaskInfo) error

// A Tasker represents an implementation of a task
// It provides the capability to Run the task given information from the previous
// task including the Target (file targeted) and Buf (output) in TaskInfo
// The run function should update the TaskInfo.Target to reflect what the Task worked on
// It should also clear and fill the Buf if the task had output
// If run returns an error the workflow will stop
// If the workflow should continue, handle the error internally including logging to Terr and return nil
type Tasker interface {
	Run(info *TaskInfo) (err error)
}

type task struct {
	Transform Transformer
	RunFunc   Runner
}

// NewTask returns a Task that will, when executed from a Workflow, transform Target with Transform(TaskInfo.Target)
// and execute RunFunc
func NewTask(t Transformer, r Runner) Tasker {
	return &task{t, r}
}

func (t *task) Run(i *TaskInfo) (err error) {
	i.Target = t.Transform(i.Src)
	return t.RunFunc(i)
}
