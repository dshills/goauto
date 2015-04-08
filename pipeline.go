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

	"gopkg.in/fsnotify.v1"
)

const (
	batchTick = 300 * time.Millisecond // batching Duration in ms
	dirOps    = fsnotify.Create | fsnotify.Rename
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
	watcher    *fsnotify.Watcher
	events     chan []fsnotify.Event
	done       chan bool
	recDirs    map[string]bool
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
	p.Watches = append(p.Watches, d)
	if p.watcher != nil {
		err = p.watcher.Add(d)
	}

	return
}

// WatchRecursive adds a GOPATH relative or absolute path to watch recursivly
func (p *Pipeline) WatchRecursive(watchDir string, ignoreHidden bool) error {
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

	// Create a watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(p.Werr, err)
		return
	}
	p.watcher = watcher

	// setup the com channels
	p.done = make(chan bool)
	p.bufferEvents()
	qdc := p.queryRecDir()
	qwc := p.queryWorkflow()

	// Add the watch directories to the watcher
	for _, w := range p.Watches {
		if err := watcher.Add(w); err != nil {
			fmt.Fprint(p.Wout, err)

		} else if p.Verbose {
			fmt.Fprintf(p.Wout, "> Watching %v\n", w)
		}
	}

	// block
	p.distributeEvents(qdc, qwc)
}

// bufferEvents watches for file events and batches them up based on a timer
// if the event distributer is busy it just keeps batching up events
// **Thanks to github.com/egonelbre for the suggestions and examples for batch events
func (p *Pipeline) bufferEvents() {
	p.events = make(chan []fsnotify.Event)

	go func() {
		defer func() {
			close(p.events)
		}()

		tick := time.Tick(batchTick)
		buf := make([]fsnotify.Event, 0, 10)
		var out chan []fsnotify.Event

		for {
			select {
			// buffer the events
			case event := <-p.watcher.Events:
				buf = append(buf, event)
			// check if we have any events
			case <-tick:
				if len(buf) > 0 {
					out = p.events
				}
			// if nil skip, otherwise send when it's ready
			case out <- buf:
				buf = make([]fsnotify.Event, 0, 10)
				out = nil
			case <-p.done:
				return
			}
		}
	}()
}

// distributeEvents sends batched events to a list of write channels
// when finished it closes the write channels
func (p *Pipeline) distributeEvents(cs ...chan<- fsnotify.Event) {
	defer func() {
		for _, c := range cs {
			close(c)
		}
	}()

	for {
		select {
		case d := <-p.events:
			if d == nil {
				return
			}
			for _, e := range d {
				for _, c := range cs {
					select {
					case c <- e:
					case <-p.done:
						return
					}
				}
			}
		}
	}
}

// queryWorkflow checks for file match for each workflow and if matches executes the workflow tasks
// returns a write channel that the caller should close
func (p *Pipeline) queryWorkflow() chan<- fsnotify.Event {
	in := make(chan fsnotify.Event)

	go func() {
		for {
			select {
			case e := <-in:
				for _, wf := range p.Workflows {
					if wf.Match(e.Name, uint32(e.Op)) {
						wf.Run(&TaskInfo{Src: e.Name, Tout: p.Wout, Terr: p.Werr, Verbose: p.Verbose})
					}
				}
			case <-p.done:
				return
			}
		}
	}()
	return in
}

// matchNewRec checks if an event is adding or renaming a directory in a recursive watch
// reruns WatchRecursive if it is
func (p *Pipeline) matchNewRec(e fsnotify.Event) {
	fi, err := os.Stat(e.Name)
	if err == nil && fi.IsDir() && dirOps&e.Op == e.Op {
		h := IsHidden(e.Name)
		for dir, iHidden := range p.recDirs {
			if h && iHidden {
				continue
			}
			if _, err := filepath.Rel(dir, e.Name); err == nil {
				if err := p.WatchRecursive(dir, iHidden); err != nil {
					fmt.Fprint(p.Wout, err)
				} else if p.Verbose {
					fmt.Fprintf(p.Wout, "> Detected new watch %v\n", e.Name)
				}
				break
			}
		}
	}
}

// queryRecDir checks if an event is adding or renaming a directory in a recursive watch
// returns a write channel that the caller should close
func (p *Pipeline) queryRecDir() chan<- fsnotify.Event {
	in := make(chan fsnotify.Event, 10) // bursts of events often come in, try not to slow the workflows down

	go func() {
		for {
			select {
			case e := <-in:
				p.matchNewRec(e)
			case <-p.done:
				return
			}
		}
	}()
	return in
}

// Stop will discontinue watching for file changes
func (p *Pipeline) Stop() (err error) {
	if p.done == nil || p.watcher == nil {
		return errors.New("Pipeline was not started or has not completed")
	}
	err = p.watcher.Close()
	close(p.done)
	if p.Verbose {
		fmt.Fprintln(p.Wout, "> Pipeline stopped")
	}
	return
}
