// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-fsnotify/fsevents"
)

type watchOSX struct {
	out         io.Writer
	eventStream *fsevents.EventStream
	done        chan struct{}
}

// NewWatchOSX returns a OSX specific file system watcher
func NewWatchOSX() Watcher {
	w := new(watchOSX)
	w.eventStream = &fsevents.EventStream{
		Paths: []string{},
		Flags: fsevents.FileEvents | fsevents.WatchRoot,
	}
	return w
}

func (w *watchOSX) SetVerbose(out io.Writer) {
	w.out = out
}

func (w *watchOSX) convertFlags(e fsevents.Event) Op {
	var f Op
	if e.Flags&fsevents.ItemCreated == fsevents.ItemCreated {
		f |= Create
	}
	if e.Flags&fsevents.ItemRemoved == fsevents.ItemRemoved {
		f |= Remove
	}
	if e.Flags&fsevents.ItemRenamed == fsevents.ItemRenamed {
		f |= Rename
	}
	if e.Flags&fsevents.ItemModified == fsevents.ItemModified {
		f |= Write
	}
	if e.Flags&fsevents.ItemInodeMetaMod == fsevents.ItemInodeMetaMod {
		f |= Chmod
	}
	return f
}

func (w *watchOSX) Start(latency time.Duration, paths []string) (<-chan ESlice, error) {
	w.done = make(chan struct{})
	c := make(chan ESlice)
	w.eventStream.Paths = paths
	w.eventStream.Latency = latency

	if w.out != nil {
		for _, d := range paths {
			fmt.Fprintln(w.out, "Watching", d)
		}
	}

	go func() {
		defer func() {
			close(c)
			if w.out != nil {
				fmt.Fprintln(w.out, "Closing watcher channel")
			}
		}()
		for {
			select {
			case msg := <-w.eventStream.Events:
				buf := make([]*Event, 0, len(msg))
				for _, e := range msg {
					buf = append(buf, &Event{Path: e.Path, Op: w.convertFlags(e)})
					if w.out != nil {
						fmt.Fprintln(w.out, Event{Path: e.Path, Op: w.convertFlags(e)})
					}
				}
				c <- buf
			case <-w.done:
				return
			}
		}
	}()

	w.eventStream.Start()
	return c, nil
}

func (w *watchOSX) Stop() error {
	if w.done == nil || w.eventStream == nil {
		return errors.New("Watcher not started or already stopped")
	}
	if w.out != nil {
		fmt.Fprintln(w.out, "Watcher stopped")
	}
	close(w.done)
	w.eventStream.Stop()
	return nil
}

func (w *watchOSX) Add(path string) error {
	w.eventStream.Paths = append(w.eventStream.Paths, path)
	return nil
}

func (w *watchOSX) Remove(path string) error {
	return nil
}
