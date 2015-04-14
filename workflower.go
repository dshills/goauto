// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

// A Workflower represents a workflow that executes a list of Taskers
type Workflower interface {
	WatchPattern(patterns ...string) error
	WatchOp(op Op)
	Add(tasks ...Tasker)
	Match(fpath string, op Op) bool
	Run(*TaskInfo)
}
