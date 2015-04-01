// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"fmt"
	"regexp"
	"time"
)

// Op describes a set of file operations.
// Mimics fsnotify
type Op uint32

// These are the generalized file operations that can trigger a notification.
// Mimics fsnotify
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// A Workflow represents a set of tasks for files matching one or more regex patterns
type Workflow struct {
	Name       string
	Concurrent bool
	Op         Op
	Regexs     []*regexp.Regexp
	Tasks      []Tasker
}

// NewWorkflow returns a Workflow with one pattern and one task
// An invlid regexp pattern will cause a panic
func NewWorkflow(name, pattern string, task Tasker) Workflower {
	return &Workflow{
		Name:       name,
		Concurrent: false,
		Op:         Create | Write | Remove | Rename,
		Regexs:     []*regexp.Regexp{regexp.MustCompile(pattern)},
		Tasks:      []Tasker{task},
	}
}

// AddPattern adds a regex for matching files for this workflow
// An invalid regexp pattern will return an error
func (wf *Workflow) AddPattern(pattern string) error {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	wf.Regexs = append(wf.Regexs, r)
	return nil
}

// Match checks a file name against the regexp of the Workflow and the file operation
func (wf *Workflow) Match(fpath string, op uint32) bool {
	if uint32(wf.Op)&op == op {
		for _, r := range wf.Regexs {
			if r.MatchString(fpath) {
				return true
			}
		}
	}
	return false
}

// AddTask adds a task to the workflow
func (wf *Workflow) AddTask(t Tasker) {
	wf.Tasks = append(wf.Tasks, t)
}

func (wf *Workflow) runner(info *TaskInfo) {
	if Verbose {
		fmt.Fprintf(info.Tout, ">> %v %v for %v\n\n", time.Now().Format("2006/01/02 3:04pm"), wf.Name, info.Src)
	}
	fname := info.Src
	info.Collect = []string{fname}
	var err error
	for _, t := range wf.Tasks {
		info.Target = "" // reset the Target
		if err = t.Run(info); err != nil {
			fmt.Fprintln(info.Terr, err)
			fmt.Fprintf(info.Terr, "Fail! Workflow did not complete for %v\n", fname)
			return
		}
		if info.Target != "" {
			// if the task set a target use it for the Src in the next task
			info.Src = info.Target
			info.Collect = append(info.Collect, info.Target)
		}
	}
}

// Run will start the execution of tasks
func (wf *Workflow) Run(info *TaskInfo) {
	if wf.Concurrent {
		go wf.runner(info)
		return
	}
	wf.runner(info)
}
