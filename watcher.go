// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"io"
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

// Event represents a file system notification
type Event struct {
	Path string
	Op   Op
}

// ESlice is an Event buffer
type ESlice []*Event

// A Watcher represents a gneric type of file system monitor
type Watcher interface {
	SetVerbose(out io.Writer)
	Start(latency time.Duration, paths []string) (<-chan ESlice, error)
	Stop() error
	Add(path string) error
	Remove(path string) error
}
