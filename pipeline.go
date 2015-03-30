// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

// A Pipeline watches one or more directories for changes
type Pipeline struct {
	Name       string
	Watches    []string
	Wout, Werr io.Writer
	Workflows  []*Workflow
	watcher    *fsnotify.Watcher
}

// NewPipeline returns a basic Pipeline with a dir to watch, output and error writers and a workflow
func NewPipeline(name string, watchDir string, wout, werr io.Writer, wf *Workflow) *Pipeline {
	p := Pipeline{Name: name, Wout: wout, Werr: werr, Workflows: []*Workflow{wf}}
	_, err := p.AddWatch(watchDir)
	if err != nil {
		panic(err)
	}
	return &p
}

// AddWatch adds a GOPATH relative or absolute path to watch
// rejects invalid paths and ignores duplicates
func (p *Pipeline) AddWatch(watchDir string) (string, error) {
	d, err := AbsPath(watchDir)
	if err != nil {
		if Verbose {
			log.Println(err)
		}
		return "", err
	}
	// Make sure we are not already watching it
	for _, w := range p.Watches {
		if w == d {
			return d, nil
		}
	}
	p.Watches = append(p.Watches, d)
	if p.watcher != nil {
		p.watcher.Add(d)
	}
	return d, nil
}

// AddRecWatch adds a GOPATH relative or absolute path to watch recursivly
func (p *Pipeline) AddRecWatch(watchDir string, ignoreHidden bool) error {
	d, err := AbsPath(watchDir)
	if err != nil {
		return err
	}
	filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			// HACKY skip hidden dir
			if (info.Name()[:1] == ".") && ignoreHidden {
				return filepath.SkipDir
			}
			p.AddWatch(path)
		}
		return nil
	})
	return nil
}

// AddWorkflow adds a workflow to the pipeline
func (p *Pipeline) AddWorkflow(w *Workflow) {
	p.Workflows = append(p.Workflows, w)
}

// Watch begins watching for changes to files in the Watches directories
// Detected file changes will be compared with workflow regexp and if match will run the workflow tasks
func (p *Pipeline) Watch(done <-chan bool) {
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintln(p.Werr, err)
		return
	}
	defer watcher.Close()
	p.watcher = watcher

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					p.doWorkflow(event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					p.doWorkflow(event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					p.doWorkflow(event.Name)
				}
			case err := <-watcher.Errors:
				fmt.Fprintln(p.Werr, "Error:", err)
			}
		}
	}()

	for _, w := range p.Watches {
		watcher.Add(w)
		if Verbose {
			log.Println("Watching", w)
		}
	}

	<-done
}

// doWorkflow checks for file match for each workflow and if matches executes the workflow tasks
func (p *Pipeline) doWorkflow(fpath string) {
	if Verbose {
		log.Println("Watcher event", fpath)
	}
	f := filepath.Base(fpath)
	for _, wf := range p.Workflows {
		if wf.Match(f) {
			wf.Run(&TaskInfo{Src: fpath, Tout: p.Wout, Terr: p.Werr})
		}
	}
}
