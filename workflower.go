// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

// A Workflower represents a workflow that executes a list of Taskers
type Workflower interface {
	WatchPattern(pattern ...string) error
	Add(t ...Tasker)
	Match(fpath string, op uint32) bool
	Run(*TaskInfo)
}
