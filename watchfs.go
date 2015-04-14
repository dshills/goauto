// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"errors"
	"fmt"
	"io"
	"time"

	"gopkg.in/fsnotify.v1"
)

type watchFS struct {
	watcher *fsnotify.Watcher
	out     io.Writer
	done    chan struct{}
}

// NewWatchFS creates a new filesystem watcher
func NewWatchFS() Watcher {
	return new(watchFS)
}

func (w *watchFS) SetVerbose(out io.Writer) {
	w.out = out
}

func (w *watchFS) Start(latency time.Duration, paths []string) (<-chan ESlice, error) {
	w.done = make(chan struct{})
	c := make(chan ESlice)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	w.watcher = watcher

	go w.bufferEvents(c, latency)

	for _, d := range paths {
		if err := w.watcher.Add(d); err != nil {
			if w.out != nil {
				fmt.Fprintln(w.out, err)
			}
		}
		if w.out != nil {
			fmt.Fprintln(w.out, "Watching", d)
		}
	}
	return c, nil
}

// bufferEvents watches for file events and batches them up based on a timer
// if the event distributer is busy it just keeps batching up events
// **Thanks to github.com/egonelbre for the suggestions and examples for batch events
func (w *watchFS) bufferEvents(send chan<- ESlice, l time.Duration) {
	defer close(send)

	tick := time.Tick(l)
	buf := make(ESlice, 0, 10)
	var out chan<- ESlice

	for {
		select {
		// buffer the events
		case e := <-w.watcher.Events:
			buf = append(buf, &Event{Path: e.Name, Op: Op(e.Op)})
		case err := <-w.watcher.Errors:
			if w.out != nil {
				fmt.Fprintln(w.out, err)
			}
			return
		// check if we have any events
		case <-tick:
			if len(buf) > 0 {
				out = send
			}
		// if nil skip, otherwise send when it's ready
		case out <- buf:
			buf = make(ESlice, 0, 10)
			out = nil
		case <-w.done:
			return
		}
	}
}

func (w *watchFS) Stop() error {
	if w.done == nil || w.watcher == nil {
		return errors.New("Watcher not started or already stopped")
	}
	if w.out != nil {
		fmt.Fprintln(w.out, "Watcher stopped")
	}
	close(w.done)
	err := w.watcher.Close()
	w.watcher = nil
	return err
}

func (w *watchFS) Add(path string) (err error) {
	if w.watcher != nil {
		if w.out != nil {
			fmt.Fprintln(w.out, "Watching", path)
		}
		return w.watcher.Add(path)
	}
	return nil
}

func (w *watchFS) Remove(path string) error {
	if w.watcher != nil {
		if w.out != nil {
			fmt.Fprintln(w.out, "Removing", path)
		}
		return w.watcher.Remove(path)
	}
	return nil
}
