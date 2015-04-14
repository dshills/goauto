// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"fmt"
	"regexp"
	"time"
)

// A Workflow represents a set of tasks for files matching one or more regex patterns
type Workflow struct {
	Name       string
	Concurrent bool
	Op         Op
	Regexs     []*regexp.Regexp
	Tasks      []Tasker
}

// NewWorkflow returns a Workflow with tasks
func NewWorkflow(tasks ...Tasker) *Workflow {
	wf := new(Workflow)
	wf.WatchOp(Create | Write | Remove | Rename)
	wf.Add(tasks...)
	return wf
}

// WatchPattern adds one or more regex for matching files for this workflow
// An invalid regexp pattern will return an error
// Sets file operations to Create|Write|Remove|Rename if not set
func (wf *Workflow) WatchPattern(patterns ...string) error {
	if wf.Op == 0 {
		wf.WatchOp(Create | Write | Remove | Rename)
	}
	for _, p := range patterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return err
		}
		wf.Regexs = append(wf.Regexs, r)
	}
	return nil
}

// WatchOp sets the file operations to match
// The default is Create | Write | Remove | Rename
func (wf *Workflow) WatchOp(op Op) {
	wf.Op = op
}

// Match checks a file name against the regexp of the Workflow and the file operation
func (wf *Workflow) Match(fpath string, op Op) bool {
	if wf.Op&op == op {
		for _, r := range wf.Regexs {
			if r.MatchString(fpath) {
				return true
			}
		}
	}
	return false
}

// Add adds a task to the workflow
func (wf *Workflow) Add(tasks ...Tasker) {
	for _, t := range tasks {
		wf.Tasks = append(wf.Tasks, t)
	}
}

func (wf *Workflow) runner(info *TaskInfo) {
	if info.Verbose {
		fmt.Fprintf(info.Tout, ">> %v %v for %v\n\n", time.Now().Format("2006/01/02 3:04pm"), wf.Name, info.Src)
	}
	fname := info.Src
	info.Collect = []string{fname}
	var err error
	for _, t := range wf.Tasks {
		info.Target = "" // reset the Target
		if err = t.Run(info); err != nil {
			fmt.Fprintln(info.Terr, err)
			fmt.Fprintf(info.Terr, "Fail! Workflow %v did not complete for %v\n\n\n", wf.Name, fname)
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
