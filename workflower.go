// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

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

// A Workflower represents a workflow that executes a list of Taskers
type Workflower interface {
	WatchPattern(patterns ...string) error
	WatchOp(op Op)
	Add(tasks ...Tasker)
	Match(fpath string, op uint32) bool
	Run(*TaskInfo)
}
