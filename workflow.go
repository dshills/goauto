// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"fmt"
	"io"
	"log"
	"regexp"
	"time"
)

// A Workflow represents a set of tasks for files matching one or more regex patterns
type Workflow struct {
	Name   string
	Regexs []*regexp.Regexp
	Tasks  []*Task
}

// NewWorkflow returns a Workflow with one pattern and one task
// An invlid regexp pattern will cause a panic
func NewWorkflow(name, pattern string, task *Task) *Workflow {
	return &Workflow{name, []*regexp.Regexp{regexp.MustCompile(pattern)}, []*Task{task}}
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
func (wf *Workflow) AddTask(t *Task) {
	wf.Tasks = append(wf.Tasks, t)
}

// Run will start the execution of tasks
func (wf *Workflow) Run(filename string, wout, werr io.Writer) {
	if Verbose {
		log.Println("Running Workflow", wf.Name, filename)
	}
	fmt.Fprintln(wout, wf.Name, time.Now())
	p := &Task{FileName: filename}
	for _, t := range wf.Tasks {
		t.FileName = p.FileName
		t.Buffer = p.Buffer
		err := t.Execute(wout, werr)
		if err != nil {
			fmt.Fprintln(werr, err)
			fmt.Fprintln(werr, "Workflow did not complete for", t.FileName)
			return
		}
		p = t
	}
	fmt.Fprintln(wout, "")
	p.Buffer.Reset()
}
