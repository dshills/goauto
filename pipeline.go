// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// NOTE check out -gcflags=-m
// want things on the stack so not GCed

const (
	batchTick = 300 * time.Millisecond // batching Duration in ms
)

// Flags to WatchRecursive to include or ignore hidden directories
const (
	IgnoreHidden  = true
	IncludeHidden = false
)

// Flags for verbose output
const (
	Verbose = true
	Silent  = false
)

// A Pipeline watches one or more directories for changes
type Pipeline struct {
	Name       string
	Watches    []string
	Wout, Werr io.Writer
	Workflows  []Workflower
	Verbose    bool
	OSX        bool
	watcher    Watcher
	recDirs    map[string]bool
	events     <-chan ESlice
}

// NewPipeline returns a basic Pipeline with a dir to watch, output and error writers and a workflow
func NewPipeline(name string, verbose bool) *Pipeline {
	p := Pipeline{Name: name, Wout: os.Stdout, Werr: os.Stderr, Verbose: verbose}
	return &p
}

// Watch adds a GOPATH relative or absolute path to watch
// rejects invalid paths and ignores duplicates
func (p *Pipeline) Watch(watchDir string) (d string, err error) {
	d, err = AbsPath(watchDir)
	if err != nil {
		if p.Verbose {
			fmt.Fprintf(p.Werr, "> %v", err)
		}
		return
	}

	// Make sure we are not already watching it
	for _, w := range p.Watches {
		if w == d {
			return
		}
	}
	if p.Verbose && p.OSX {
		fmt.Fprintf(p.Wout, "OSX watches are always recursive and do not skip directories. Adding %v recursivly\n", d)
	}
	p.Watches = append(p.Watches, d)
	if p.watcher != nil {
		err = p.watcher.Add(d)
	}

	return
}

// WatchRecursive adds a GOPATH relative or absolute path to watch recursivly
func (p *Pipeline) WatchRecursive(watchDir string, ignoreHidden bool) error {
	if p.OSX {
		_, err := p.Watch(watchDir)
		return err
	}
	d, err := AbsPath(watchDir)
	if err != nil {
		return err
	}
	if p.recDirs == nil {
		p.recDirs = make(map[string]bool)
	}
	p.recDirs[d] = ignoreHidden
	err = filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if IsHidden(info.Name()) && ignoreHidden {
				return filepath.SkipDir
			}
			_, err = p.Watch(path)
		}
		return nil
	})
	return nil
}

// Add adds one or more Workflows to the pipeline
func (p *Pipeline) Add(ws ...Workflower) {
	for _, w := range ws {
		p.Workflows = append(p.Workflows, w)
	}
}

// Start begins watching for changes to files in the Watches directories
// Detected file changes will be compared with workflow regexp and if match will run the workflow tasks
func (p *Pipeline) Start() {
	if p.watcher == nil {
		if p.OSX {
			p.watcher = NewWatchOSX()
		} else {
			p.watcher = NewWatchFS()
		}
		if p.Verbose {
			p.watcher.SetVerbose(p.Wout)
		}
	}

	if p.Wout == nil {
		p.Wout = os.Stdout
	}
	if p.Werr == nil {
		p.Werr = os.Stderr
	}
	if p.Name == "" {
		p.Name = "<UNNAMED>"
	}

	if len(p.Watches) < 1 {
		fmt.Fprintln(p.Werr, "Pipeline", p.Name, "is not watching anything")
	}

	if len(p.Workflows) < 1 {
		fmt.Fprintln(p.Werr, "Pipeline", p.Name, "has no Workflows")
	}

	// setup the com channels
	qdc := p.queryRecDir()
	qwc := p.queryWorkflow()

	var err error
	p.events, err = p.watcher.Start(batchTick, p.Watches)
	if err != nil {
		fmt.Fprintln(p.Werr, err)
		return
	}

	// block
	p.distributeEvents(qdc, qwc)
}

// distributeEvents sends batched events to a list of write channels
// when finished it closes the write channels
func (p *Pipeline) distributeEvents(cs ...chan<- *Event) {
	defer func() {
		for _, c := range cs {
			close(c)
		}
	}()

	for {
		select {
		case d := <-p.events:
			if d == nil || len(d) < 1 {
				return
			}
			for _, e := range d {
				for _, c := range cs {
					c <- e
				}
			}
		}
	}
}

// queryWorkflow checks for file match for each workflow and if matches executes the workflow tasks
// returns a write channel that the caller should close
func (p *Pipeline) queryWorkflow() chan<- *Event {
	in := make(chan *Event)

	go func() {
		for {
			select {
			case e := <-in:
				if e == nil {
					return
				}
				for _, wf := range p.Workflows {
					if wf.Match(e.Path, e.Op) {
						wf.Run(&TaskInfo{Src: e.Path, Tout: p.Wout, Terr: p.Werr, Verbose: p.Verbose})
					}
				}
			}
		}
	}()
	return in
}

// matchNewRec checks if an event is adding or renaming a directory in a recursive watch
// reruns WatchRecursive if it is
func (p *Pipeline) matchNewRec(e Event) {
	dirOps := Create | Rename
	fi, err := os.Stat(e.Path)
	if err == nil && fi.IsDir() && dirOps&e.Op == e.Op {
		h := IsHidden(e.Path)
		for dir, iHidden := range p.recDirs {
			if h && iHidden {
				continue
			}
			if _, err := filepath.Rel(dir, e.Path); err == nil {
				if err := p.WatchRecursive(dir, iHidden); err != nil {
					fmt.Fprint(p.Wout, err)
				} else if p.Verbose {
					fmt.Fprintf(p.Wout, "> Detected new watch %v\n", e.Path)
				}
				break
			}
		}
	}
}

// queryRecDir checks if an event is adding or renaming a directory in a recursive watch
// returns a write channel that the caller should close
func (p *Pipeline) queryRecDir() chan<- *Event {
	in := make(chan *Event, 10) // bursts of events often come in, try not to slow the workflows down

	go func() {
		for {
			select {
			case e := <-in:
				if e == nil {
					return
				}
				p.matchNewRec(*e)
			}
		}
	}()
	return in
}

// Stop will discontinue watching for file changes
func (p *Pipeline) Stop() (err error) {
	if p.watcher == nil {
		return errors.New("Pipeline was not started or has not completed")
	}
	err = p.watcher.Stop()
	if err != nil {
		fmt.Fprintln(p.Wout, err)
	}
	if p.Verbose {
		fmt.Fprintln(p.Wout, "> Pipeline stopped")
	}
	return
}
