// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

// A Workflow represents a set of tasks for files matching one or more regex patterns
type Workflow struct {
	Name   string
	Regexs []*regexp.Regexp
	Tasks  []Tasker
}

// NewWorkflow returns a Workflow with one pattern and one task
// An invlid regexp pattern will cause a panic
func NewWorkflow(name, pattern string, task Tasker) *Workflow {
	return &Workflow{name, []*regexp.Regexp{regexp.MustCompile(pattern)}, []Tasker{task}}
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

// Match checks a file name against the regexp of the Workflow
func (wf Workflow) Match(fn string) bool {
	for _, r := range wf.Regexs {
		if r.MatchString(fn) {
			return true
		}
	}
	return false
}

// AddTask adds a task to the workflow
func (wf *Workflow) AddTask(t Tasker) {
	wf.Tasks = append(wf.Tasks, t)
}

// Run will start the execution of tasks
func (wf *Workflow) Run(info *TaskInfo) {
	if Verbose {
		log.Println("Running Workflow", wf.Name, info.Src)
	}
	fmt.Fprintln(info.Tout, wf.Name, time.Now())
	fname := info.Src

	var err error
	for _, t := range wf.Tasks {
		info.Target = "" // reset the Target
		if err = t.Run(info); err != nil {
			fmt.Fprintln(info.Terr, err)
			fmt.Fprintln(info.Terr, "Workflow did not complete for", fname)
			return
		}
		if info.Target != "" {
			// if the task set a target use it for the Src in the next task
			info.Src = info.Target
		}
	}
	fmt.Fprintln(info.Tout, "")
}
